package upload_test

import (
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"testing"
)

func TestFilter_AllowAllFiles(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", true},
		{"testdata/SampleJPGImage.jpg", true},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", true},
		{"testdata/SampleText.txt", true},
		{"testdata/SampleVideo.mp4", true},
		{"testdata/ScreenShotJPG.jpg", true},
		{"testdata/ScreenShotPNG.png", true},
	}

	t.Run("WithEmptyPatterns", func(t *testing.T) {
		f := upload.NewFilter([]string{""}, []string{""}, true)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("WithWildCardPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"*"}, []string{""}, true)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})
}

func TestFilter_AllowPNGFiles(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", false},
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", false},
		{"testdata/SampleText.txt", false},
		{"testdata/SampleVideo.mp4", false},
		{"testdata/ScreenShotJPG.jpg", false},
		{"testdata/ScreenShotPNG.png", true},
	}

	f := upload.NewFilter([]string{"*.png"}, []string{""}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_AllowPNGFilesWithFolder(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/folder/SampleJPGImage.jpg", false},
		{"testdata/folder/SamplePNGImage.png", true},

	}

	f := upload.NewFilter([]string{"*.png"}, []string{""}, false)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_AllowPNGAndJPGFiles(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", false},
		{"testdata/SampleJPGImage.jpg", true},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", false},
		{"testdata/SampleText.txt", false},
		{"testdata/SampleVideo.mp4", false},
		{"testdata/ScreenShotJPG.jpg", true},
		{"testdata/ScreenShotPNG.png", true},
	}

	f := upload.NewFilter([]string{"*.png", "*.jpg"}, []string{""}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_AllowSampleImageFiles(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", false},
		{"testdata/SampleJPGImage.jpg", true},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", true},
		{"testdata/SampleText.txt", false},
		{"testdata/SampleVideo.mp4", false},
		{"testdata/ScreenShotJPG.jpg", false},
		{"testdata/ScreenShotPNG.png", false},
	}

	f := upload.NewFilter([]string{"Sample*"}, []string{"*.mp3", "*.txt", "*.mp4"}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_DisallowAllFiles(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", false},
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", false},
		{"testdata/SampleSVGImage.svg", false},
		{"testdata/SampleText.txt", false},
		{"testdata/SampleVideo.mp4", false},
		{"testdata/ScreenShotJPG.jpg", false},
		{"testdata/ScreenShotPNG.png", false},
	}

	f := upload.NewFilter([]string{"*"}, []string{"*"}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_DisallowScreenShots(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", true},
		{"testdata/SampleJPGImage.jpg", true},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", true},
		{"testdata/SampleText.txt", true},
		{"testdata/SampleVideo.mp4", true},
		{"testdata/ScreenShotJPG.jpg", false},
		{"testdata/ScreenShotPNG.png", false},
	}

	f := upload.NewFilter([]string{""}, []string{"*ScreenShot*"}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_DisallowVideos(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleAudio.mp3", true},
		{"testdata/SampleJPGImage.jpg", true},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/SampleSVGImage.svg", true},
		{"testdata/SampleText.txt", true},
		{"testdata/SampleVideo.mp4", false},
		{"testdata/ScreenShotJPG.jpg", true},
		{"testdata/ScreenShotPNG.png", true},
	}

	f := upload.NewFilter([]string{""}, []string{""}, false)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}
}

func TestFilter_DisallowFolder(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/folder/SampleJPGImage.jpg", false},
		{"testdata/folder/SamplePNGImage.png", false},

	}

	f := upload.NewFilter([]string{"*.png"}, []string{"folder"}, false)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}
