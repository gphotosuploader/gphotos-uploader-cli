package upload_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

func TestWalker_GetAllFiles(t *testing.T) {
	var includePatterns = []string{"_ALL_FILES_"}
	var excludePatterns = []string{""}
	var expected = []string{
		"testdata/SampleAudio.mp3",
		"testdata/SampleJPGImage.jpg",
		"testdata/SamplePNGImage.png",
		"testdata/SampleSVGImage.svg",
		"testdata/SampleText.txt",
		"testdata/SampleVideo.mp4",
		"testdata/ScreenShotJPG.jpg",
		"testdata/ScreenShotPNG.png",
		"testdata/folder1/SamplePNGImage.png",
		"testdata/folder1/SampleJPGImage.jpg",
		"testdata/folder2/SamplePNGImage.png",
		"testdata/folder2/SampleJPGImage.jpg",
		"testdata/folder-symlink/SamplePNGImage.png",
		"testdata/folder-symlink/SampleJPGImage.jpg",
	}

	got, err := getIncludedFilesByScanFolder(includePatterns, excludePatterns)

	require.NoError(t, err)
	assert.ElementsMatch(t, expected, got)
}

func TestWalker_GetAllPNGFiles(t *testing.T) {
	var includePatterns = []string{"**/*.png"}
	var excludePatterns = []string{""}
	var expected = []string{
		"testdata/SamplePNGImage.png",
		"testdata/ScreenShotPNG.png",
		"testdata/folder1/SamplePNGImage.png",
		"testdata/folder2/SamplePNGImage.png",
		"testdata/folder-symlink/SamplePNGImage.png",
	}

	got, err := getIncludedFilesByScanFolder(includePatterns, excludePatterns)

	require.NoError(t, err)
	assert.ElementsMatch(t, expected, got)
}

func TestWalker_GetAllFilesExcludeFolder1(t *testing.T) {
	var includePatterns = []string{"_ALL_FILES_"}
	var excludePatterns = []string{"folder1"}
	var expected = []string{
		"testdata/SampleAudio.mp3",
		"testdata/SampleJPGImage.jpg",
		"testdata/SamplePNGImage.png",
		"testdata/SampleSVGImage.svg",
		"testdata/SampleText.txt",
		"testdata/SampleVideo.mp4",
		"testdata/ScreenShotJPG.jpg",
		"testdata/ScreenShotPNG.png",
		"testdata/folder2/SamplePNGImage.png",
		"testdata/folder2/SampleJPGImage.jpg",
		"testdata/folder-symlink/SamplePNGImage.png",
		"testdata/folder-symlink/SampleJPGImage.jpg",
	}

	got, err := getIncludedFilesByScanFolder(includePatterns, excludePatterns)

	require.NoError(t, err)
	assert.ElementsMatch(t, expected, got)
}

func TestRelativePath(t *testing.T) {
	var objectsTest = []struct {
		base string
		in   string
		want string
	}{
		{base: "/foo/bar", in: "/foo/bar/xyz", want: "xyz"},
		{base: "/foo/bar/", in: "/foo/bar/xyz", want: "xyz"},
		{base: "/foo/bar", in: "/foo/bar/xyz/", want: "xyz"},
		{base: "/foo/bar", in: "foo/bar/xyz", want: "foo/bar/xyz"},
		{base: "/foo/bar", in: "/foo/bar", want: "."},
		{base: "/foo/bar/", in: "/foo/bar", want: "."},
		{base: "/foo/bar", in: "/foo/bar/", want: "."},
		{base: "", in: "/foo/bar", want: "/foo/bar"},
		{base: "/foo/bar", in: "/abc/def", want: "/abc/def"},
	}

	for _, tc := range objectsTest {
		assert.Equal(t, tc.want, upload.RelativePath(tc.base, tc.in))
	}
}

func getIncludedFilesByScanFolder(includePatterns []string, excludePatterns []string) ([]string, error) {
	ft := &mock.FileTracker{
		MarkAsUploadedFn: func(path string) error {
			return nil
		},
		IsUploadedFn: func(path string) bool {
			return strings.Contains(path, "AlreadyUploaded")
		},
		UnmarkAsUploadedFn: func(path string) error {
			return nil
		},
	}
	filterFiles := filter.MustCompile(includePatterns, excludePatterns)

	u := upload.UploadFolderJob{
		FileTracker:  ft,
		SourceFolder: "testdata",
		CreateAlbums: "Off",
		Filter:       filterFiles,
	}

	foundItems, err := u.ScanFolder(&mock.Logger{})
	if err != nil {
		return nil, err
	}

	var results []string
	for _, i := range foundItems {
		results = append(results, i.Path)
	}

	return results, nil
}
