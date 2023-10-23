package upload

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		assert.Equal(t, tc.want, f.Name())
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

		assert.Equal(t, tc.want, f.String())
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
	err := appFS.MkdirAll("src/", 0755)
	require.NoError(t, err)

	err = afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644)
	require.NoError(t, err)

	for _, tc := range testCases {
		f := NewFileItem(tc.in)
		_, size, err := f.Open()

		if tc.errExpected {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.wantSize, size)
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
	err := appFS.MkdirAll("src/", 0755)
	require.NoError(t, err)

	err = afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644)
	require.NoError(t, err)

	for _, tc := range testCases {
		f := NewFileItem(tc.in)

		assert.Equal(t, tc.want, f.Size())
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
	err := appFS.MkdirAll("src/", 0755)
	require.NoError(t, err)

	err = afero.WriteFile(appFS, "src/existent", []byte("this is content of existing file"), 0644)
	require.NoError(t, err)

	for _, tc := range testCases {
		f := NewFileItem(tc.in)
		err := f.Remove()

		if tc.errExpected {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestFileItem_GroupByAlbum(t *testing.T) {
	items := []FileItem{
		{Path: "file1.jpg", AlbumName: "album 1"},
		{Path: "file2.jpg", AlbumName: "album 2"},
		{Path: "file3.jpg", AlbumName: "album 1"},
		{Path: "file4.jpg", AlbumName: "album 2"},
		{Path: "file5.jpg", AlbumName: "album 3"},
	}

	expectedGroups := map[string][]FileItem{
		"album 1": {
			{Path: "file1.jpg", AlbumName: "album 1"},
			{Path: "file3.jpg", AlbumName: "album 1"},
		},
		"album 2": {
			{Path: "file2.jpg", AlbumName: "album 2"},
			{Path: "file4.jpg", AlbumName: "album 2"},
		},
		"album 3": {
			{Path: "file5.jpg", AlbumName: "album 3"},
		},
	}

	groupedItems := GroupByAlbum(items)

	assert.Len(t, groupedItems, len(expectedGroups))

	for albumName, expectedItems := range expectedGroups {

		assert.Contains(t, groupedItems, albumName)
		assert.Equal(t, expectedItems, groupedItems[albumName])
	}
}
