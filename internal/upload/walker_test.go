package upload_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

func TestWalker_GetAllFiles(t *testing.T) {
	var includePatterns = []string{""}
	var excludePatterns = []string{""}
	var allowVideos = true

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":            true,
		"testdata/SampleJPGImage.jpg":         true,
		"testdata/SamplePNGImage.png":         true,
		"testdata/SampleSVGImage.svg":         true,
		"testdata/SampleText.txt":             true,
		"testdata/SampleVideo.mp4":            true,
		"testdata/ScreenShotJPG.jpg":          true,
		"testdata/ScreenShotPNG.png":          true,
		"testdata/folder1/SamplePNGImage.png": true,
		"testdata/folder1/SampleJPGImage.jpg": true,
		"testdata/folder2/SamplePNGImage.png": true,
		"testdata/folder2/SampleJPGImage.jpg": true,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns, allowVideos)
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
	var includePatterns = []string{"*.png"}
	var excludePatterns = []string{""}
	var allowVideos = false

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":            false,
		"testdata/SampleJPGImage.jpg":         false,
		"testdata/SamplePNGImage.png":         true,
		"testdata/SampleSVGImage.svg":         false,
		"testdata/SampleText.txt":             false,
		"testdata/SampleVideo.mp4":            false,
		"testdata/ScreenShotJPG.jpg":          false,
		"testdata/ScreenShotPNG.png":          true,
		"testdata/folder1/SamplePNGImage.png": true,
		"testdata/folder1/SampleJPGImage.jpg": false,
		"testdata/folder2/SamplePNGImage.png": true,
		"testdata/folder2/SampleJPGImage.jpg": false,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns, allowVideos)
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
	var includePatterns = []string{""}
	var excludePatterns = []string{"folder1"}
	var allowVideos = true

	var want = map[string]bool{
		"testdata/SampleAudio.mp3":            true,
		"testdata/SampleJPGImage.jpg":         true,
		"testdata/SamplePNGImage.png":         true,
		"testdata/SampleSVGImage.svg":         true,
		"testdata/SampleText.txt":             true,
		"testdata/SampleVideo.mp4":            true,
		"testdata/ScreenShotJPG.jpg":          true,
		"testdata/ScreenShotPNG.png":          true,
		"testdata/folder1/SamplePNGImage.png": false,
		"testdata/folder1/SampleJPGImage.jpg": false,
		"testdata/folder2/SamplePNGImage.png": true,
		"testdata/folder2/SampleJPGImage.jpg": true,
	}

	got, err := getScanFolderResult(includePatterns, excludePatterns, allowVideos)
	if err != nil {
		t.Fatalf("no error was expected at this point: err=%s", err)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("want: %v, got: %v, file: %s", want[i], got[i], i)
		}
	}
}

func getScanFolderResult(includePatterns []string, excludePatterns []string, allowVideos bool) (map[string]bool, error) {
	ft := &mock.FileTracker{
		CacheAsAlreadyUploadedFn: func(path string) error {
			return nil
		},
		IsAlreadyUploadedFn: func(path string) (bool, error) {
			return false, nil
		},
		RemoveAsAlreadyUploadedFn: func(path string) error {
			return nil
		},
	}
	u := upload.UploadFolderJob{
		FileTracker:        ft,
		SourceFolder:       "testdata",
		CreateAlbum:        false,
		CreateAlbumBasedOn: "",
		Filter:             upload.NewFilter(includePatterns, excludePatterns, allowVideos),
	}

	foundItems, err := u.ScanFolder(&mock.Logger{})
	if err != nil {
		return nil, err
	}

	var results = map[string]bool{}
	for _, i := range foundItems {
		results[i.Path] = true
	}

	return results, nil
}
