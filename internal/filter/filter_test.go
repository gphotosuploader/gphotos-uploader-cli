package filter_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
)

func TestCompile(t *testing.T) {
	testCases := []struct {
		name         string
		allowedList  []string
		excludedList []string
		errExpected  bool
	}{
		{name: "empty patterns", allowedList: []string{""}, excludedList: []string{""}, errExpected: false},
		{name: "valid patterns", allowedList: []string{"**"}, excludedList: []string{"**/*.png"}, errExpected: false},
		{name: "invalid allowed list", allowedList: []string{"[]a]"}, excludedList: []string{""}, errExpected: true},
		{name: "invalid excluded list", allowedList: []string{""}, excludedList: []string{"[]a]"}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := filter.Compile(tc.allowedList, tc.excludedList)
			if tc.errExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMustCompile(t *testing.T) {
	testCases := []struct {
		name          string
		allowedList   []string
		excludedList  []string
		panicExpected bool
	}{
		{name: "empty patterns", allowedList: []string{""}, excludedList: []string{""}, panicExpected: false},
		{name: "valid patterns", allowedList: []string{"**"}, excludedList: []string{"**/*.png"}, panicExpected: false},
		{name: "invalid allowed list", allowedList: []string{"[]a]"}, excludedList: []string{""}, panicExpected: true},
		{name: "invalid excluded list", allowedList: []string{""}, excludedList: []string{"[]a]"}, panicExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.panicExpected {
						t.Error("panic was not expected but the function panic")
					}
				}
			}()

			_ = filter.MustCompile(tc.allowedList, tc.excludedList)
			if tc.panicExpected {
				t.Error("panic was expected but the function doesn't panic")
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
		{"testdata/SampleJPGImage.JPG", true},
		{"testdata/SamplePNGImage.PNG", true},
		{"testdata/SampleSVGImage.SVG", false},
	}

	t.Run("ByUsingEmptyPatterns", func(t *testing.T) {
		f, err := filter.Compile([]string{""}, []string{""})

		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
		}
	})

	t.Run("ByUsingRepeatedEmptyPatterns", func(t *testing.T) {
		f, err := filter.Compile([]string{"", "", ""}, []string{"", "", ""})

		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.Compile([]string{"_IMAGE_EXTENSIONS_"}, []string{""})

		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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
		{"testdata/SampleJPGImage.JPG", true},
		{"testdata/SamplePNGImage.PNG", true},
		{"testdata/SampleSVGImage.SVG", true},
	}

	t.Run("ByUsingWildCardPattern", func(t *testing.T) {
		f, err := filter.Compile([]string{"**"}, []string{""})

		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.Compile([]string{"_ALL_FILES_"}, []string{""})

		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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

	f, err := filter.Compile([]string{"**/*.png"}, []string{""})

	require.NoError(t, err)

	for _, tc := range testCases {
		assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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

	f, err := filter.Compile([]string{"**/*.png", "**/*.jpg"}, []string{""})
	require.NoError(t, err)

	for _, tc := range testCases {
		assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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

	f, err := filter.Compile([]string{"**/Sample*"}, []string{"**/*.mp3", "**/*.txt", "**/*.mp4"})
	require.NoError(t, err)

	for _, tc := range testCases {
		assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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
		f, err := filter.Compile([]string{"**"}, []string{"**"})
		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
		}
	})

	t.Run("ByUsingTaggedPattern", func(t *testing.T) {
		f, err := filter.Compile([]string{"_ALL_FILES_"}, []string{"_ALL_FILES_"})
		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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

	f, err := filter.Compile([]string{"_ALL_FILES_"}, []string{"**/ScreenShot*"})
	require.NoError(t, err)

	for _, tc := range testCases {
		assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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
		f, err := filter.Compile([]string{"_ALL_FILES_"}, []string{"_ALL_VIDEO_FILES_"})
		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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

	f, err := filter.Compile([]string{"**/*.png"}, []string{"*/folder1/*"})
	require.NoError(t, err)

	for _, tc := range testCases {
		assert.Equal(t, tc.out, f.IsAllowed(tc.file))
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
		f, err := filter.Compile([]string{""}, []string{"**/folder1/*"})
		require.NoError(t, err)

		for _, tc := range testCases {
			assert.Equal(t, tc.out, f.IsExcluded(tc.file))
		}
	})

}
