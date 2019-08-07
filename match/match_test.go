package match_test

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/match"
)

var matchTests = []struct {
	pattern string
	path    string
	match   bool
}{
	{"", "", true},
	{"", "foo", true},
	{"", "/x/y/z/foo", true},
	{"*.go", "/foo/bar/test.go", true},
	{"*.c", "/foo/bar/test.go", false},
	{"*", "/foo/bar/test.go", true},
	{"foo*", "/foo/bar/test.go", true},
	{"bar*", "/foo/bar/test.go", true},
	{"/bar*", "/foo/bar/test.go", false},
	{"bar/*", "/foo/bar/test.go", true},
	{"baz/*", "/foo/bar/test.go", false},
	{"bar/test.go", "/foo/bar/test.go", true},
	{"bar/*.go", "/foo/bar/test.go", true},
	{"ba*/*.go", "/foo/bar/test.go", true},
	{"bb*/*.go", "/foo/bar/test.go", false},
	{"test.*", "/foo/bar/test.go", true},
	{"tesT.*", "/foo/bar/test.go", false},
	{"bar/*", "/foo/bar/baz", true},
	{"bar", "/foo/bar", true},
	{"/foo/bar", "/foo/bar", true},
	{"/foo/bar/", "/foo/bar", true},
	{"/foo/bar", "/foo/baz", false},
	{"/foo/bar", "/foo/baz/", false},
	{"/foo///bar", "/foo/bar", true},
	{"/foo/../bar", "/foo/bar", false},
	{"/foo/../bar", "/bar", true},
	{"/foo", "/foo/baz", true},
	{"/foo/", "/foo/baz", true},
	{"/foo/*", "/foo", false},
	{"/foo/*", "/foo/baz", true},
	{"bar", "/foo/bar/baz", true},
	{"bar", "/foo/bar/test.go", true},
	{"/foo/*test.*", "/foo/bar/test.go", false},
	{"/foo/*/test.*", "/foo/bar/test.go", true},
	{"/foo/*/bar/test.*", "/foo/bar/test.go", false},
	{"/*/*/bar/test.*", "/foo/bar/test.go", false},
	{"/*/*/bar/test.*", "/foo/bar/baz/test.go", false},
	{"/*/*/baz/test.*", "/foo/bar/baz/test.go", true},
	{"/*/foo/bar/test.*", "/foo/bar/baz/test.go", false},
	{"/*/foo/bar/test.*", "/foo/bar/baz/test.go", false},
	{"/foo/bar/test.*", "bar/baz/test.go", false},
	{"/x/y/bar/baz/test.*", "bar/baz/test.go", false},
	{"/x/y/bar/baz/test.c", "bar/baz/test.go", false},
	{"baz/test.*", "bar/baz/test.go", true},
	{"baz/tesT.*", "bar/baz/test.go", false},
	{"test.go", "bar/baz/test.go", true},
	{"*.go", "bar/baz/test.go", true},
	{"*.c", "bar/baz/test.go", false},
	{"sdk", "/foo/bar/sdk", true},
	{"sdk", "/foo/bar/sdk/test/sdk_foo.go", true},
	{
		"sdk/*/cpp/*/*vars*.html",
		"/usr/share/doc/libreoffice/sdk/docs/cpp/ref/a00517.html",
		false,
	},
	{"foo/**/bar/*.go", "/home/user/foo/work/special/project/bar/test.go", true},
	{"foo/**/bar/*.go", "/home/user/foo/bar/test.go", true},
	{"foo/**/bar/*.go", "x/foo/bar/test.go", true},
	{"foo/**/bar/*.go", "foo/bar/test.go", true},
	{"foo/**/bar/*.go", "foo/bar/baz/bar/test.go", true},
	{"foo/**/bar/*.go", "/home/user/foo/test.c", false},
	{"foo/**/bar/*.go", "bar/foo/main.go", false},
	{"foo/**/bar/*.go", "/foo/bar/main.go", true},
	{"foo/**/bar/*.go", "bar/main.go", false},
	{"foo/**/bar", "/home/user/foo/x/y/bar", true},
	{"foo/**/bar", "/home/user/foo/x/y/bar/main.go", true},
	{"foo/**/bar/**/x", "/home/user/foo/bar/x", true},
	{"foo/**/bar/**/x", "/home/user/foo/blaaa/blaz/bar/shared/work/x", true},
	{"user/**/important*", "/home/user/work/x/y/hidden/x", false},
	{"user/**/hidden*/**/c", "/home/user/work/x/y/hidden/z/a/b/c", true},
	{"c:/foo/*test.*", "c:/foo/bar/test.go", false},
	{"c:/foo", "c:/foo/bar", true},
	{"c:/foo/", "c:/foo/bar", true},
	{"c:/foo/*/test.*", "c:/foo/bar/test.go", true},
	{"c:/foo/*/bar/test.*", "c:/foo/bar/test.go", false},
}

func testPattern(t *testing.T, pattern, path string, shouldMatch bool) {
	match, err := match.Match(pattern, path)
	if err != nil {
		t.Errorf("test pattern %q failed: expected no error for path %q, but error returned: %v",
			pattern, path, err)
	}

	if match != shouldMatch {
		t.Errorf("test: filter.Match(%q, %q): expected %v, got %v",
			pattern, path, shouldMatch, match)
	}
}

func TestMatch(t *testing.T) {
	for _, test := range matchTests {
		t.Run("", func(t *testing.T) {
			testPattern(t, test.pattern, test.path, test.match)
		})

		// Test with native path separator
		if filepath.Separator != '/' {
			pattern := strings.Replace(test.pattern, "/", string(filepath.Separator), -1)
			// Test with pattern as native
			t.Run("pattern-native", func(t *testing.T) {
				testPattern(t, pattern, test.path, test.match)
			})

			path := strings.Replace(test.path, "/", string(filepath.Separator), -1)
			t.Run("path-native", func(t *testing.T) {
				// Test with path as native
				testPattern(t, test.pattern, path, test.match)
			})

			t.Run("both-native", func(t *testing.T) {
				// Test with both pattern and path as native
				testPattern(t, pattern, path, test.match)
			})
		}
	}

	t.Run("TestMatchWithEmptyString", func(t *testing.T) {
		_, err := match.Match("*.go", "")
		if err != match.ErrBadString {
			t.Errorf("test empty string failed: expected %s error, but different error returned: %v", match.ErrBadString, err)
		}
	})
}

func ExampleMatch() {
	match, _ := match.Match("*.go", "/home/user/file.go")
	fmt.Printf("match: %v\n", match)
	// Output:
	// match: true
}

func ExampleMatch_wildcards() {
	match, _ := match.Match("/home/[uU]ser/?.go", "/home/user/F.go")
	fmt.Printf("match: %v\n", match)
	// Output:
	// match: true
}

var filterMatchOneTests = []struct {
	patterns []string
	path     string
	match    bool
}{
	{[]string{"*.go"}, "/foo/bar/test.go", true},
	{[]string{"*.c"}, "/foo/bar/test.go", false},
	{[]string{"*.go", "*.c"}, "/foo/bar/test.go", true},
	{[]string{"*"}, "/foo/bar/test.go", true},
	{[]string{"x"}, "/foo/bar/test.go", false},
	{[]string{"?"}, "/foo/bar/test.go", false},
	{[]string{"?", "x"}, "/foo/bar/x", true},
	{[]string{"/*/*/bar/test.*"}, "/foo/bar/test.go", false},
	{[]string{"/*/*/bar/test.*", "*.go"}, "/foo/bar/test.go", true},
	{[]string{"", "*.c"}, "/foo/bar/test.go", false},
}

func TestMatchOne(t *testing.T) {
	for i, test := range filterMatchOneTests {
		match, err := match.MatchOne(test.patterns, test.path)
		if err != nil {
			t.Errorf("test %d failed: expected no error for patterns %q, but error returned: %v",
				i, test.patterns, err)
			continue
		}

		if match != test.match {
			t.Errorf("test %d: filter.MatchOne(%q, %q): expected %v, got %v",
				i, test.patterns, test.path, test.match, match)
		}
	}

	t.Run("TestMatchOneWithEmptyString", func(t *testing.T) {
		_, err := match.MatchOne([]string{"*.go"}, "")
		if err != match.ErrBadString {
			t.Errorf("test empty string failed: expected %s error, but different error returned: %v", match.ErrBadString, err)
		}
	})
}

func ExampleMatchOne() {
	match, _ := match.MatchOne([]string{"*.c", "*.go"}, "/home/user/file.go")
	fmt.Printf("match: %v\n", match)
	// Output:
	// match: true
}

func extractTestLines(t testing.TB) (lines []string) {
	f, err := os.Open("testdata/libreoffice.txt.bz2")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	sc := bufio.NewScanner(bzip2.NewReader(f))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines
}

func TestFilterPatternsFile(t *testing.T) {
	lines := extractTestLines(t)

	var testPatterns = []struct {
		pattern string
		hits    uint
	}{
		{"*.html", 18249},
		{"sdk", 22186},
		{"sdk/*/cpp/*/*vars.html", 3},
	}

	for _, test := range testPatterns {
		var c uint
		for _, line := range lines {
			match, err := match.Match(test.pattern, line)
			if err != nil {
				t.Error(err)
				continue
			}

			if match {
				c++
				// fmt.Printf("pattern %q, line %q\n", test.pattern, line)
			}
		}

		if c != test.hits {
			t.Errorf("wrong number of hits for pattern %q: want %d, got %d",
				test.pattern, test.hits, c)
		}
	}
}
