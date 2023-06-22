package upload_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

func TestWalker_GetAllFiles(t *testing.T) {
	var includePatterns = []string{"_ALL_FILES_"}
	var excludePatterns = []string{""}

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":                   true,
		"testdata/SampleJPGImage.jpg":                true,
		"testdata/SamplePNGImage.png":                true,
		"testdata/SampleSVGImage.svg":                true,
		"testdata/SampleText.txt":                    true,
		"testdata/SampleVideo.mp4":                   true,
		"testdata/ScreenShotJPG.jpg":                 true,
		"testdata/ScreenShotPNG.png":                 true,
		"testdata/folder1/SamplePNGImage.png":        true,
		"testdata/folder1/SampleJPGImage.jpg":        true,
		"testdata/folder2/SamplePNGImage.png":        true,
		"testdata/folder2/SampleJPGImage.jpg":        true,
		"testdata/folder-symlink/SamplePNGImage.png": true,
		"testdata/folder-symlink/SampleJPGImage.jpg": true,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns)
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("want: %v, got: %v, file: %s", want[i], got[i], i)
		}
	}
}

func TestWalker_GetAllPNGFiles(t *testing.T) {
	var includePatterns = []string{"**/*.png"}
	var excludePatterns = []string{""}

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":                   false,
		"testdata/SampleJPGImage.jpg":                false,
		"testdata/SamplePNGImage.png":                true,
		"testdata/SampleSVGImage.svg":                false,
		"testdata/SampleText.txt":                    false,
		"testdata/SampleVideo.mp4":                   false,
		"testdata/ScreenShotJPG.jpg":                 false,
		"testdata/ScreenShotPNG.png":                 true,
		"testdata/folder1/SamplePNGImage.png":        true,
		"testdata/folder1/SampleJPGImage.jpg":        false,
		"testdata/folder2/SamplePNGImage.png":        true,
		"testdata/folder2/SampleJPGImage.jpg":        false,
		"testdata/folder-symlink/SamplePNGImage.png": true,
		"testdata/folder-symlink/SampleJPGImage.jpg": false,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns)
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("want: %v, got: %v, file: %s", want[i], got[i], i)
		}
	}
}

func TestWalker_GetAllFilesExcludeFolder1(t *testing.T) {
	var includePatterns = []string{"_ALL_FILES_"}
	var excludePatterns = []string{"folder1"}

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":                   true,
		"testdata/SampleJPGImage.jpg":                true,
		"testdata/SamplePNGImage.png":                true,
		"testdata/SampleSVGImage.svg":                true,
		"testdata/SampleText.txt":                    true,
		"testdata/SampleVideo.mp4":                   true,
		"testdata/ScreenShotJPG.jpg":                 true,
		"testdata/ScreenShotPNG.png":                 true,
		"testdata/folder1/SamplePNGImage.png":        false,
		"testdata/folder1/SampleJPGImage.jpg":        false,
		"testdata/folder2/SamplePNGImage.png":        true,
		"testdata/folder2/SampleJPGImage.jpg":        true,
		"testdata/folder-symlink/SamplePNGImage.png": true,
		"testdata/folder-symlink/SampleJPGImage.jpg": true,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns)
	if err != nil {
		t.Fatalf("no error was expected at this point: err=%s", err)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("want: %v, got: %v, file: %s", want[i], got[i], i)
		}
	}
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

		got := upload.RelativePath(tc.base, tc.in)
		if got != tc.want {
			t.Errorf("Test Case (%s), basepath '%s': want '%s', got '%s'", tc.base, tc.in, tc.want, got)
		}
	}
}

func getScanFolderResult(includePatterns []string, excludePatterns []string) (map[string]bool, error) {
	var results = map[string]bool{}
	ft := &mock.FileTracker{
		MarkAsUploadedFn: func(path string) error {
			return nil
		},
		IsUploadedFn: func(path string) bool {
			return false
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

	for _, i := range foundItems {
		results[i.Path] = true
	}

	return results, nil
}
