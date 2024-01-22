package upload

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlbumName(t *testing.T) {
	var testData = []struct {
		name  string
		album string

		in   string
		want string
	}{
		{
			name:  "album set an album's name",
			album: "name:albumName",
			in:    "/foo/bar/file.jpg",
			want:  "albumName",
		},
		{
			name:  "album set an album's name based on folder path",
			album: "auto:folderPath",
			in:    "/foo/bar/file.jpg",
			want:  "foo_bar",
		},
		{
			name:  "album set an album's name based on folder name",
			album: "auto:folderName",
			in:    "/foo/bar/file.jpg",
			want:  "bar",
		},
		{
			name:  "album set an album's name with unexpected key (not `name` or `auto`)",
			album: "foo:bar",
			in:    "/foo/bar/file.jpg",
			want:  "",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job := UploadFolderJob{
				Album: tt.album,
			}

			assert.Equal(t, tt.want, job.albumName(tt.in, ""))
		})

	}
}

func TestAlbumNameWithInvalidParameter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("A Panic was expected but not reached.")
		}
	}()
	job := UploadFolderJob{
		Album: "auto:fooBar",
	}
	_ = job.albumName("/foo/bar/file.jpg", "")
}

func TestAlbumNameUsingFolderPath(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{in: "", out: ""},
		{in: "foo", out: ""},
		{in: "foo/", out: "foo"},
		{in: "foo/bar", out: "foo"},
		{in: "foo/bar/", out: "foo_bar"},
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "foo_bar"},
	}
	for _, tt := range testData {
		assert.Equal(t, tt.out, albumNameUsingFolderPath(tt.in))
	}
}

func TestAlbumNameUsingFolderName(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{in: "", out: ""},
		{in: "foo", out: ""},
		{in: "foo/", out: "foo"},
		{in: "foo/bar", out: "foo"},
		{in: "foo/bar/", out: "bar"},
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "bar"},
	}
	for _, tt := range testData {
		assert.Equal(t, tt.out, albumNameUsingFolderName(tt.in))
	}
}

func TestGetTemplateFunctionName(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{in: "", out: ""},
		{in: "$cutLeft(", out: "cutLeft"},
		{in: "$cutLeft (", out: ""},
		{in: "cutLeft(", out: ""},
		{in: "$cutLeft(anything", out: "cutLeft"},
	}

	for _, tt := range testData {
		assert.Equal(t, tt.out, getTemplateFunctionName(tt.in))
	}
}

func TestParseAlbumNameTampleWithTokens(t *testing.T) {
	timeObj := time.Date(2034, time.December, 31, 16, 5, 59, 0, time.UTC)
	filePath := "/foo/bar/file.jpg"

	var testData = []struct {
		in  string
		out string
	}{
		{in: "%_year%", out: "2034"},
		{in: "%_day%", out: "31"},
		{in: "%_month%", out: "12"},

		{in: "%_folderpath%", out: "foo_bar"},
		{in: "%_directory%", out: "bar"},
		{in: "%_parent_directory%", out: "foo"},

		{in: "%_time%", out: "16:05:59"},
		{in: "%_time_en%", out: "04:05:59 PM"},
		{in: "%_year%_%_year%", out: "2034_2034"},
	}

	for _, tt := range testData {
		val, err := parseAlbumNameTemplate(tt.in, filePath, timeObj)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, tt.out, val)
	}
}

func TestParseAlbumNameTampleWithFunctions(t *testing.T) {
	timeObj := time.Date(2034, time.December, 31, 16, 5, 59, 0, time.UTC)
	filePath := "/foo/bar/file.jpg"

	var testData = []struct {
		in  string
		out string
	}{
		{in: "$cutLeft(%_year%, 2)", out: "34"},
		{in: "$cutLeft(%_year%, 5)", out: ""},
		{in: "$cutRight(%_year%, 2)", out: "20"},
		{in: "$cutRight(%_year%, 5)", out: ""},

		{in: "$lower(Hello World!)", out: "hello world!"},
		{in: "$upper(Hello World!)", out: "HELLO WORLD!"},
		{in: "$sentence(Hello World!)", out: "Hello world!"},
		{in: "$sentence(Hello World!)", out: "Hello world!"},
		{in: "$title(Title of a Set of Photos)", out: "Title Of A Set Of Photos"},

		{in: "$regexp(Hello World!, World, Universe)", out: "Hello Universe!"},
		{in: "$regexp(Hello _World!,_, ',')", out: "Hello ,World!"},
		{in: "$regexp(Hello _World!,'[_!]', ',')", out: "Hello ,World,"},
	}

	for _, tt := range testData {
		val, err := parseAlbumNameTemplate(tt.in, filePath, timeObj)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.Equal(t, tt.out, val)
	}
}

func TestParseAlbumNameTampleWithInvalidParameter(t *testing.T) {
	// testing  for Y2k bug ;)
	timeObj := time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC)

	var testData = []struct {
		in  string
		err string
	}{
		{in: "%_ABC%", err: "invalid token: ABC"},
		{in: "$ABC(Z)", err: "unknown function: ABC"},
		{in: "$cutLeft(Z,Z)", err: "cutLeft requires a number as second argument"},
		{in: "$cutLeft(Z,Z, Z)", err: "cutLeft requires 2 arguments"},
		{in: "$cutLeft(Z)", err: "cutLeft requires 2 arguments"},
		{in: "$cutLeft($cutLeft(Z)", err: "function missing closing parenthesis"},
		{in: "$cutLeft($cutLeft(Z), 2)", err: "cutLeft requires 2 arguments"},
		{in: "$cutLeft($cutRight(Z), 2)", err: "cutRight requires 2 arguments"},
		{in: "$lower()", err: "lower requires 1 argument"},
		{in: "$regexp(Hello World!, ^[a-z+\\[$, Universe)", err: "invalid regexp pattern: ^[a-z+\\[$"},
		{in: "$regexp(Hello World!, _, ABC'()')", err: "Can't mix quoted & unquoted content in function arg: regexp"},
		{in: "$regexp(Hello World!, _, ')", err: "string missing closing quote"},
	}

	for _, tt := range testData {
		_, err := parseAlbumNameTemplate(tt.in, "", timeObj)
		assert.EqualError(t, err, tt.err)
	}
}
