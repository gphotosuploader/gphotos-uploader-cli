package upload_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/upload"
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

	t.Run("ByUsingEmptyPatterns", func(t *testing.T) {
		f := upload.NewFilter([]string{""}, []string{""}, true)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingRepeatedEmptyPatterns", func(t *testing.T) {
		f := upload.NewFilter([]string{"", "", ""}, []string{"", "", ""}, true)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingWildCardPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"*"}, []string{""}, true)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"_ALL_FILES_"}, []string{""}, true)
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
		{"testdata/folder/SampleJPGImage.jpg", true},
		{"testdata/folder/SamplePNGImage.png", true},
	}

	f := upload.NewFilter([]string{"*.png", "*.jpg"}, []string{""}, false)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_AllowImageFilesStartingWithSample(t *testing.T) {
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

	f := upload.NewFilter([]string{"Sample*"}, []string{"*.mp3", "*.txt", "*.mp4"}, false)
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

	t.Run("ByUsingWildcardPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"*"}, []string{"*"}, false)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"_ALL_FILES_"}, []string{"_ALL_FILES_"}, false)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

}

func TestFilter_DisallowFilesStartingWithScreenShot(t *testing.T) {
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

	t.Run("ByUsingParameter", func(t *testing.T) {
		f := upload.NewFilter([]string{"*"}, []string{""}, false)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f := upload.NewFilter([]string{"_ALL_FILES_"}, []string{"_ALL_VIDEO_FILES_"}, false)
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

}

func TestFilter_DisallowAFolder(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/folder/SampleJPGImage.jpg", false},
		{"testdata/folder/SamplePNGImage.png", false},
	}

	f := upload.NewFilter([]string{"*.png"}, []string{"folder"}, true)
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}
