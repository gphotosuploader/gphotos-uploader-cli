package filter_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
)

func TestFilter_Validate(t *testing.T) {
	testCases := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		errExpected     bool
	}{
		{name: "empty patterns", includePatterns: []string{""}, excludePatterns: []string{""}, errExpected: false},
		{name: "valid patterns", includePatterns: []string{"**"}, excludePatterns: []string{"**/*.png"}, errExpected: false},
		{name: "invalid includePattern", includePatterns: []string{"[]a]"}, excludePatterns: []string{""}, errExpected: true},
		{name: "invalid excludePattern", includePatterns: []string{""}, excludePatterns: []string{"[]a]"}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := filter.New(tc.includePatterns, tc.excludePatterns)
			if err != nil && !tc.errExpected {
				t.Errorf("error was not expected, got: %v", err)
			}
		})
	}
}

func TestFilter_AllowDefaultFiles(t *testing.T) {
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

	t.Run("ByUsingEmptyPatterns", func(t *testing.T) {
		f, err := filter.New([]string{""}, []string{""})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingRepeatedEmptyPatterns", func(t *testing.T) {
		f, err := filter.New([]string{"", "", ""}, []string{"", "", ""})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.New([]string{"_IMAGE_EXTENSIONS_"}, []string{""})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})
}

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

	t.Run("ByUsingWildCardPattern", func(t *testing.T) {
		f, err := filter.New([]string{"**"}, []string{""})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.New([]string{"_ALL_FILES_"}, []string{""})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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

	f, err := filter.New([]string{"**/*.png"}, []string{""})
	if err != nil {
		t.Fatalf("error was not expected at this point: %v", err)
	}
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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

	f, err := filter.New([]string{"**/*.png", "**/*.jpg"}, []string{""})
	if err != nil {
		t.Fatalf("error was not expected at this point: %v", err)
	}
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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

	f, err := filter.New([]string{"**/Sample*"}, []string{"**/*.mp3", "**/*.txt", "**/*.mp4"})
	if err != nil {
		t.Fatalf("error was not expected at this point: %v", err)
	}
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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
		f, err := filter.New([]string{"**"}, []string{"**"})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.New([]string{"_ALL_FILES_"}, []string{"_ALL_FILES_"})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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

	f, err := filter.New([]string{"_ALL_FILES_"}, []string{"**/ScreenShot*"})
	if err != nil {
		t.Fatalf("error was not expected at this point: %v", err)
	}
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
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

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.New([]string{"_ALL_FILES_"}, []string{"_ALL_VIDEO_FILES_"})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsAllowed(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

}

func TestFilter_IncludingPNGExceptAFolder(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", true},
		{"testdata/folder1/SampleJPGImage.jpg", false},
		{"testdata/folder1/SamplePNGImage.png", false},
		{"testdata/folder2/SampleJPGImage.jpg", false},
		{"testdata/folder2/SamplePNGImage.png", true},
	}

	f, err := filter.New([]string{"**/*.png"}, []string{"*/folder1/*"})
	if err != nil {
		t.Fatalf("error was not expected at this point: %v", err)
	}
	for _, tc := range testCases {
		got := f.IsAllowed(tc.file)
		if tc.out != got {
			t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
		}
	}

}

func TestFilter_ExcludingAFolder(t *testing.T) {
	var testCases = []struct {
		file string
		out  bool
	}{
		{"testdata/SampleJPGImage.jpg", false},
		{"testdata/SamplePNGImage.png", false},
		{"testdata/folder1/SampleJPGImage.jpg", true},
		{"testdata/folder1/SamplePNGImage.png", true},
		{"testdata/folder2/SampleJPGImage.jpg", false},
		{"testdata/folder2/SamplePNGImage.png", false},
	}

	t.Run("ExcludingFolder1", func(t *testing.T) {
		f, err := filter.New([]string{""}, []string{"**/folder1/*"})
		if err != nil {
			t.Fatalf("error was not expected at this point: %v", err)
		}
		for _, tc := range testCases {
			got := f.IsExcluded(tc.file)
			if tc.out != got {
				t.Errorf("Filter result was not expected: file=%s, want %t, got %t", tc.file, tc.out, got)
			}
		}
	})

}
