package upload_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

func TestFileItem_Open(t *testing.T) {
	var testCases = []struct {
		name        string
		in          string
		wantSize    int64
		errExpected bool
	}{
		{name: "ShouldReturnErrorWhenFileDoesNotExist", in: "testdata/non-existent.jpg", wantSize: 0, errExpected: true},
		{name: "ShouldReturnSuccessWhenFileExists", in: "testdata/SampleJPGImage.jpg", wantSize: 51085, errExpected: false},
	}

	for _, tc := range testCases {
		f := upload.FileItem(tc.in)
		_, size, err := f.Open()
		switch {
		case tc.errExpected && err == nil:
			t.Errorf("TestCase(%s), error was expected, but not happened", tc.name)
		case !tc.errExpected && err != nil:
			t.Errorf("TestCase(%s), error was not expected: err=%s", tc.name, err)
		case size != tc.wantSize:
			t.Errorf("TestCase(%s), want: %d, got: %d", tc.name, tc.wantSize, size)
		}
	}
}

func TestFileItem_Name(t *testing.T) {
	var testCases = []struct {
		in   string
		want string
	}{
		{in: "testdata/SampleJPGImage.jpg", want: "SampleJPGImage.jpg"},
		{in: "testdata/SamplePNGImage.png", want: "SamplePNGImage.png"},
		{in: "testdata/SampleSVGImage.svg", want: "SampleSVGImage.svg"},
	}

	for _, tc := range testCases {
		f := upload.FileItem(tc.in)
		if got := f.Name(); got != tc.want {
			t.Errorf("TestCase(%s), want: %s, got: %s", tc.in, tc.want, got)
		}
	}
}

func TestFileItem_String(t *testing.T) {
	var testCases = []struct {
		in   string
		want string
	}{
		{in: "testdata/SampleJPGImage.jpg", want: "testdata/SampleJPGImage.jpg"},
		{in: "testdata/SamplePNGImage.png", want: "testdata/SamplePNGImage.png"},
		{in: "testdata/SampleSVGImage.svg", want: "testdata/SampleSVGImage.svg"},
	}

	for _, tc := range testCases {
		f := upload.FileItem(tc.in)
		if got := f.String(); got != tc.want {
			t.Errorf("TestCase(%s), want: %s, got: %s", tc.in, tc.want, got)
		}
	}
}

func TestFileItem_Size(t *testing.T) {
	var testCases = []struct {
		in   string
		want int64
	}{
		{in: "testdata/SampleJPGImage.jpg", want: 51085},
		{in: "testdata/SamplePNGImage.png", want: 104327},
		{in: "testdata/SampleSVGImage.svg", want: 24276},
	}

	for _, tc := range testCases {
		f := upload.FileItem(tc.in)
		if got := f.Size(); got != tc.want {
			t.Errorf("Test Case(%s), want: %d, got: %d", tc.in, tc.want, got)
		}
	}
}
