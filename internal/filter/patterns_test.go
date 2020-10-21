package filter

import (
	"reflect"
	"testing"
)

func Test_TranslatePatterns(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  []string
	}{
		{name: "_ALL_FILES_", input: []string{"_ALL_FILES_"}, want: []string{"**"}},
		{name: "_ALL_VIDEO_FILES_", input: []string{"_ALL_VIDEO_FILES_"}, want: []string{"**/*.mpg", "**/*.mod", "**/*.mmv", "**/*.tod", "**/*.wmv", "**/*.asf", "**/*.avi", "**/*.divx", "**/*.mov", "**/*.m4v", "**/*.3gp", "**/*.3g2", "**/*.mp4", "**/*.m2t", "**/*.m2ts", "**/*.mts", "**/*.mkv"}},
		{name: "Without Tagged Pattern", input: []string{"**/*.png"}, want: []string{"**/*.png"}},
		{name: "Tagged & Not Tagged Patterns", input: []string{"**/*.png", "_ALL_VIDEO_FILES_"}, want: []string{"**/*.png", "**/*.mpg", "**/*.mod", "**/*.mmv", "**/*.tod", "**/*.wmv", "**/*.asf", "**/*.avi", "**/*.divx", "**/*.mov", "**/*.m4v", "**/*.3gp", "**/*.3g2", "**/*.mp4", "**/*.m2t", "**/*.m2ts", "**/*.mts", "**/*.mkv"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := translatePatterns(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}
}

func Test_DeleteEmpty(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  []string
	}{
		{name: "one element", input: []string{"foo"}, want: []string{"foo"}},
		{name: "three elements", input: []string{"foo", "bar", "baz"}, want: []string{"foo", "bar", "baz"}},
		{name: "elements w/ empty one", input: []string{"foo", "", "bar"}, want: []string{"foo", "bar"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := deleteEmpty(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want: %#v, got: %#v", tc.want, got)
			}
		})
	}

	t.Run("empty input array", func(t *testing.T) {
		got := deleteEmpty([]string{})
		if nil != got {
			t.Errorf("want: []string{nil}, got: %#v", got)
		}
	})
}

func Test_ValidatePatterns(t *testing.T) {
	testCases := []struct {
		name        string
		input       []string
		errExpected bool
	}{
		{name: "valid pattern returns nil", input: []string{"**"}, errExpected: false},
		{name: "two valid patterns returns nil", input: []string{"**/*.png", "**/*.jpg"}, errExpected: false},
		{name: "invalid pattern returns error", input: []string{"[]a]"}, errExpected: true},
		{name: "valid and invalid patterns returns error", input: []string{"**", "[]a]"}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := validatePatterns(tc.input)
			if got != nil && !tc.errExpected {
				t.Errorf("error was not expected, got: %v", got)
			}
		})
	}
}

func Test_MatchAnyPattern(t *testing.T) {
	testCases := []struct {
		name        string
		patterns    []string
		input       string
		shouldMatch bool
		errExpected bool
	}{
		{name: "input match pattern", patterns: []string{"*"}, input: "foo", shouldMatch: true, errExpected: false},
		{name: "input doesn't match pattern", patterns: []string{"*"}, input: "foo/bar.jpg", shouldMatch: false, errExpected: false},
		{name: "invalid pattern returns error", patterns: []string{"[]a]"}, input: "]", shouldMatch: false, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := matchAnyPattern(tc.patterns, tc.input)
			if tc.shouldMatch != got || (err != nil && !tc.errExpected) {
				t.Errorf("want: %v, got: %v, err: %v", tc.shouldMatch, got, err)
			}
		})
	}
}
