package push

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
)

func TestGetOrCreateAlbum(t *testing.T) {
	errCreate := errors.New("create failed")

	testCases := []struct {
		name         string
		title        string
		existing     map[string]string
		createErr    error
		wantID       string
		wantErr      error
		wantCreated  []string
		wantAPICalls bool
	}{
		{name: "Should return empty ID without calling the API when title is empty", title: "", wantID: ""},
		{name: "Should return the existing album", title: "Holidays", existing: map[string]string{"Holidays": "id-holidays"}, wantID: "id-holidays", wantAPICalls: true},
		{name: "Should create the album when it doesn't exist", title: "New", wantID: "id-New", wantCreated: []string{"New"}, wantAPICalls: true},
		{name: "Should fail when album creation fails", title: "New", createErr: errCreate, wantErr: errCreate, wantAPICalls: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &fakeAlbums{byTitle: tc.existing, createErr: tc.createErr}

			got, err := getOrCreateAlbum(context.Background(), svc, tc.title)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantID, got)
			assert.Equal(t, tc.wantCreated, svc.created)
			assert.Equal(t, tc.wantAPICalls, svc.calls > 0)
		})
	}
}

func TestNewPhotosService(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		photos, err := newPhotosService(http.DefaultClient, &fakeUploadSessionTracker{}, &mock.Logger{})
		require.NoError(t, err)
		assert.NotNil(t, photos.Uploader)
		assert.NotNil(t, photos.Albums)
		assert.NotNil(t, photos.MediaItems)
	})

	t.Run("Should fail when HTTP client is nil", func(t *testing.T) {
		_, err := newPhotosService(nil, &fakeUploadSessionTracker{}, &mock.Logger{})
		assert.Error(t, err)
	})
}

func TestPushCmd_push(t *testing.T) {
	t.Run("Should upload all files without album", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{}, "a.jpg", "b.jpg")

		require.NoError(t, env.run(false))

		assert.ElementsMatch(t, env.paths("a.jpg", "b.jpg"), env.media.createdFiles())
		assert.Equal(t, []string{""}, env.media.albumIDs())
		assert.ElementsMatch(t, env.paths("a.jpg", "b.jpg"), env.tracker.uploadedPaths())
		assert.Zero(t, env.albums.calls)
		assert.FileExists(t, env.path("a.jpg"))
	})

	t.Run("Should skip already uploaded and not allowed files", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{}, "new.jpg", "old.jpg", "notes.txt")
		env.tracker.uploaded[env.path("old.jpg")] = true

		require.NoError(t, env.run(false))

		assert.Equal(t, env.paths("new.jpg"), env.media.createdFiles())
	})

	t.Run("Should create the album when it doesn't exist", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{Album: "name:Holidays"}, "a.jpg")

		require.NoError(t, env.run(false))

		assert.Equal(t, []string{"Holidays"}, env.albums.created)
		assert.Equal(t, []string{"id-Holidays"}, env.media.albumIDs())
	})

	t.Run("Should reuse an existing album", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{Album: "name:Holidays"}, "a.jpg")
		env.albums.byTitle = map[string]string{"Holidays": "existing-id"}

		require.NoError(t, env.run(false))

		assert.Empty(t, env.albums.created)
		assert.Equal(t, []string{"existing-id"}, env.media.albumIDs())
	})

	t.Run("Should skip files when the album can't be created", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{Album: "name:Holidays"}, "a.jpg")
		env.albums.createErr = errors.New("create failed")

		require.NoError(t, env.run(false))

		assert.Empty(t, env.media.createdFiles())
		assert.Empty(t, env.tracker.uploadedPaths())
		assert.True(t, env.logger.FailfInvoked)
	})

	t.Run("Should not upload anything in dry run mode", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{DeleteAfterUpload: true}, "a.jpg")

		require.NoError(t, env.run(true))

		assert.Zero(t, env.uploader.calls)
		assert.Empty(t, env.tracker.uploadedPaths())
		assert.FileExists(t, env.path("a.jpg"))
		assert.True(t, env.logger.InfoInvoked)
	})

	t.Run("Should stop when daily quota is exceeded", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{}, "a.jpg", "b.jpg")
		env.uploader.err = func(string) error { return &gphotos.ErrDailyQuotaExceeded{} }

		err := env.run(false)

		var quotaErr *gphotos.ErrDailyQuotaExceeded
		assert.ErrorAs(t, err, &quotaErr)
		assert.Equal(t, 1, env.uploader.calls)
		assert.Empty(t, env.tracker.uploadedPaths())
	})

	t.Run("Should continue when a file fails to upload", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{DeleteAfterUpload: true}, "bad.jpg", "good.jpg")
		bad := env.path("bad.jpg")
		env.uploader.err = func(path string) error {
			if path == bad {
				return errors.New("upload failed")
			}
			return nil
		}

		require.NoError(t, env.run(false))

		assert.Equal(t, env.paths("good.jpg"), env.tracker.uploadedPaths())
		assert.FileExists(t, bad, "failed uploads must not be deleted")
		assert.NoFileExists(t, env.path("good.jpg"))
		assert.True(t, env.logger.FailfInvoked)
	})

	t.Run("Should delete files after upload when configured", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{DeleteAfterUpload: true}, "a.jpg")

		require.NoError(t, env.run(false))

		assert.NoFileExists(t, env.path("a.jpg"))
	})

	t.Run("Should log when file deletion fails", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{DeleteAfterUpload: true}, "a.jpg")
		// Remove the file while it's being uploaded, so the deletion fails.
		env.uploader.err = func(path string) error { return os.Remove(path) }

		require.NoError(t, env.run(false))

		assert.True(t, env.logger.ErrorfInvoked)
	})

	t.Run("Should continue when the file can't be tracked as uploaded", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{}, "a.jpg", "b.jpg")
		env.tracker.markErr = errors.New("tracker failed")

		require.NoError(t, env.run(false))

		assert.Len(t, env.media.createdFiles(), 2)
		assert.True(t, env.logger.WarnfInvoked)
	})

	t.Run("Should fail when filter patterns are invalid", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{IncludePatterns: []string{"["}}, "a.jpg")

		err := env.run(false)

		assert.ErrorContains(t, err, "include patterns are invalid")
		assert.Zero(t, env.uploader.calls)
	})

	t.Run("Should process every configured job", func(t *testing.T) {
		env := newPushEnv(t, config.FolderUploadJob{}, "a.jpg")
		other := t.TempDir()
		writeFiles(t, other, "c.jpg")
		env.cli.Config.Jobs = append(env.cli.Config.Jobs, config.FolderUploadJob{SourceFolder: other})

		require.NoError(t, env.run(false))

		assert.ElementsMatch(t, []string{env.path("a.jpg"), filepath.Join(other, "c.jpg")}, env.media.createdFiles())
	})
}

// pushEnv holds a push command wired to fake services and a temporary source folder.
type pushEnv struct {
	t        *testing.T
	dir      string
	cli      *app.App
	logger   *mock.Logger
	tracker  *fakeFileTracker
	uploader *fakeUploader
	albums   *fakeAlbums
	media    *fakeMediaItems
}

func newPushEnv(t *testing.T, job config.FolderUploadJob, files ...string) *pushEnv {
	t.Helper()
	dir := t.TempDir()
	writeFiles(t, dir, files...)
	job.SourceFolder = dir

	logger := &mock.Logger{}
	tracker := &fakeFileTracker{uploaded: map[string]bool{}}
	return &pushEnv{
		t:   t,
		dir: dir,
		cli: &app.App{
			FileTracker: tracker,
			Logger:      logger,
			Config:      &config.Config{Jobs: []config.FolderUploadJob{job}},
		},
		logger:   logger,
		tracker:  tracker,
		uploader: &fakeUploader{},
		albums:   &fakeAlbums{},
		media:    &fakeMediaItems{},
	}
}

func (e *pushEnv) run(dryRun bool) error {
	cmd := &PushCmd{
		GlobalFlags: &flags.GlobalFlags{Debug: true}, // Debug hides the progress bar.
		DryRunMode:  dryRun,
	}
	photos := &gphotos.Client{Uploader: e.uploader, Albums: e.albums, MediaItems: e.media}
	return cmd.push(context.Background(), e.cli, photos)
}

func (e *pushEnv) path(name string) string {
	return filepath.Join(e.dir, name)
}

func (e *pushEnv) paths(names ...string) []string {
	result := make([]string, 0, len(names))
	for _, n := range names {
		result = append(result, e.path(n))
	}
	return result
}

func writeFiles(t *testing.T, dir string, names ...string) {
	t.Helper()
	for _, n := range names {
		require.NoError(t, os.WriteFile(filepath.Join(dir, n), []byte("content of "+n), 0600))
	}
}

// fakeFileTracker implements app.FileTracker in memory.
type fakeFileTracker struct {
	uploaded map[string]bool
	marked   []string
	markErr  error
}

func (f *fakeFileTracker) MarkAsUploaded(path string) error {
	if f.markErr != nil {
		return f.markErr
	}
	f.marked = append(f.marked, path)
	return nil
}
func (f *fakeFileTracker) IsUploaded(path string) bool   { return f.uploaded[path] }
func (f *fakeFileTracker) UnmarkAsUploaded(string) error { return nil }
func (f *fakeFileTracker) Close() error                  { return nil }
func (f *fakeFileTracker) Destroy() error                { return nil }

func (f *fakeFileTracker) uploadedPaths() []string {
	sort.Strings(f.marked)
	return f.marked
}

// fakeUploadSessionTracker implements app.UploadSessionTracker.
type fakeUploadSessionTracker struct{}

func (f *fakeUploadSessionTracker) Get(string) (string, bool) { return "", false }
func (f *fakeUploadSessionTracker) Set(string, string)        {}
func (f *fakeUploadSessionTracker) Delete(string)             {}
func (f *fakeUploadSessionTracker) Close()                    {}

// fakeUploader implements gphotos.MediaUploader.
type fakeUploader struct {
	calls int
	err   func(path string) error
}

func (f *fakeUploader) UploadFile(_ context.Context, path string) (string, error) {
	f.calls++
	if f.err != nil {
		if err := f.err(path); err != nil {
			return "", err
		}
	}
	return "token-" + path, nil
}

// fakeAlbums implements gphotos.AlbumsService.
type fakeAlbums struct {
	byTitle   map[string]string
	createErr error
	created   []string
	calls     int
}

func (f *fakeAlbums) GetByTitle(_ context.Context, title string) (*albums.Album, error) {
	f.calls++
	if id, ok := f.byTitle[title]; ok {
		return &albums.Album{ID: id, Title: title}, nil
	}
	return nil, errors.New("album not found")
}

func (f *fakeAlbums) Create(_ context.Context, title string) (*albums.Album, error) {
	f.calls++
	if f.createErr != nil {
		return nil, f.createErr
	}
	f.created = append(f.created, title)
	return &albums.Album{ID: "id-" + title, Title: title}, nil
}

func (f *fakeAlbums) AddMediaItems(context.Context, string, []string) error { return nil }
func (f *fakeAlbums) GetById(context.Context, string) (*albums.Album, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeAlbums) List(context.Context) ([]albums.Album, error) { return nil, nil }
func (f *fakeAlbums) PaginatedList(context.Context, *albums.PaginatedListOptions) ([]albums.Album, string, error) {
	return nil, "", nil
}

// fakeMediaItems implements gphotos.MediaItemsService, recording created items.
type fakeMediaItems struct {
	created []createdItem
}

type createdItem struct {
	albumID  string
	filename string
}

func (f *fakeMediaItems) CreateToAlbum(_ context.Context, albumID string, item media_items.SimpleMediaItem) (*media_items.MediaItem, error) {
	f.created = append(f.created, createdItem{albumID: albumID, filename: item.Filename})
	return &media_items.MediaItem{ID: "media-" + item.Filename, Filename: item.Filename}, nil
}

func (f *fakeMediaItems) createdFiles() []string {
	var files []string
	for _, c := range f.created {
		files = append(files, c.filename)
	}
	sort.Strings(files)
	return files
}

// albumIDs returns the distinct album IDs used to create media items.
func (f *fakeMediaItems) albumIDs() []string {
	seen := map[string]bool{}
	var ids []string
	for _, c := range f.created {
		if !seen[c.albumID] {
			seen[c.albumID] = true
			ids = append(ids, c.albumID)
		}
	}
	return ids
}

func (f *fakeMediaItems) Create(context.Context, media_items.SimpleMediaItem) (*media_items.MediaItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMediaItems) CreateMany(context.Context, []media_items.SimpleMediaItem) ([]*media_items.MediaItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMediaItems) CreateManyToAlbum(context.Context, string, []media_items.SimpleMediaItem) ([]*media_items.MediaItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMediaItems) Get(context.Context, string) (*media_items.MediaItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMediaItems) ListByAlbum(context.Context, string) ([]*media_items.MediaItem, error) {
	return nil, nil
}
func (f *fakeMediaItems) PaginatedList(context.Context, *media_items.PaginatedListOptions) ([]media_items.MediaItem, string, error) {
	return nil, "", nil
}
