package upload

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

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
		f := NewFileItem(tc.in)
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
		f := NewFileItem(tc.in)
		if got := f.String(); got != tc.want {
			t.Errorf("TestCase(%s), want: %s, got: %s", tc.in, tc.want, got)
		}
	}
}

func TestFileItem_Open(t *testing.T) {
	var testCases = []struct {
		name        string
		in          string
		wantSize    int64
		errExpected bool
	}{
		{name: "ShouldReturnErrorWhenFileDoesNotExist", in: "src/non-existent", wantSize: 0, errExpected: true},
		{name: "ShouldReturnSuccessWhenFileExists", in: "src/existent", wantSize: 32, errExpected: false},
	}

	appFS = afero.NewMemMapFs()
	// create test files and directories
	if err := appFS.MkdirAll("src/", 0755); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}
	if err := afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}

	for _, tc := range testCases {
		f := NewFileItem(tc.in)
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

func TestFileItem_Size(t *testing.T) {
	var testCases = []struct {
		name string
		in   string
		want int64
	}{
		{name: "ShouldReturnZeroWhenFileDoesNotExist", in: "src/non-existent", want: 0},
		{name: "ShouldReturnSizeWhenFileExists", in: "src/existent", want: 32},
	}

	appFS = afero.NewMemMapFs()
	// create test files and directories
	if err := appFS.MkdirAll("src/", 0755); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}
	if err := afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}

	for _, tc := range testCases {
		f := NewFileItem(tc.in)
		if got := f.Size(); got != tc.want {
			t.Errorf("Test Case(%s), want: %d, got: %d", tc.in, tc.want, got)
		}
	}
}

func TestFileItem_Remove(t *testing.T) {
	var testCases = []struct {
		name        string
		in          string
		errExpected bool
	}{
		{name: "ShouldErrorWhenFileDoesNotExist", in: "src/non-existent", errExpected: true},
		{name: "ShouldReturnSuccessWhenFileExists", in: "src/existent", errExpected: false},
	}

	appFS = afero.NewMemMapFs()
	// create test files and directories
	if err := appFS.MkdirAll("src/", 0755); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}
	if err := afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644); err != nil {
		t.Fatalf("error was not expected at this point: err=%s", err)
	}

	for _, tc := range testCases {
		f := NewFileItem(tc.in)
		err := f.Remove()
		switch {
		case tc.errExpected && err == nil:
			t.Errorf("TestCase(%s), error was expected, but not happened", tc.name)
		case !tc.errExpected && err != nil:
			t.Errorf("TestCase(%s), error was not expected: err=%s", tc.name, err)
		}
	}
}

func TestFileItem_GroupByAlbum(t *testing.T) {
	items := []FileItem{
		{Path: "file1.jpg", AlbumName: "Album 1"},
		{Path: "file2.jpg", AlbumName: "Album 2"},
		{Path: "file3.jpg", AlbumName: "Album 1"},
		{Path: "file4.jpg", AlbumName: "Album 2"},
		{Path: "file5.jpg", AlbumName: "Album 3"},
	}

	expectedGroups := map[string][]FileItem{
		"Album 1": {
			{Path: "file1.jpg", AlbumName: "Album 1"},
			{Path: "file3.jpg", AlbumName: "Album 1"},
		},
		"Album 2": {
			{Path: "file2.jpg", AlbumName: "Album 2"},
			{Path: "file4.jpg", AlbumName: "Album 2"},
		},
		"Album 3": {
			{Path: "file5.jpg", AlbumName: "Album 3"},
		},
	}

	groupedItems := GroupByAlbum(items)

	if len(groupedItems) != len(expectedGroups) {
		t.Errorf("Unexpected number of groups. Expected: %d, Got: %d", len(expectedGroups), len(groupedItems))
	}

	for albumName, expectedItems := range expectedGroups {
		actualItems, ok := groupedItems[albumName]
		if !ok {
			t.Errorf("Expected group '%s' not found", albumName)
			continue
		}

		if !reflect.DeepEqual(actualItems, expectedItems) {
			t.Errorf("Mismatched items for group '%s'. Expected: %v, Got: %v", albumName, expectedItems, actualItems)
		}
	}
}
