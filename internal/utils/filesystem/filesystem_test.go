package filesystem_test

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/Flaque/filet"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/utils/filesystem"
)

func TestAbsolutePath(t *testing.T) {

	t.Run("WithAbsolutePaths", func(t *testing.T) {
		var absolutePathInputs = []struct {
			in  string
			out string
		}{
			{"/", "/"},
			{"/xyz", "/xyz"},
			{"/xyz/./abc", "/xyz/abc"},
			{"/xyz/../abc", "/abc"},
			{"/xyz/abc/..", "/xyz"},
			{"/xyz/../abc/..", "/"},
			{"/xyz/../..", "/"},
			{"/xyz///abc/..", "/xyz"},
		}

		for _, test := range absolutePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			if got != test.out {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, test.out, got)
			}
		}
	})

	t.Run("WithRelativePath", func(t *testing.T) {
		var relativePathInputs = []struct {
			in  string
			out string
		}{
			{"", ""},
			{"./", ""},
			{"xyz", "xyz"},
			{"xyz/./abc", "xyz/abc"},
			{"xyz/../abc", "abc"},
			{"xyz/abc/..", "xyz"},
			{"xyz/../abc/..", ""},
			{"xyz/../..", ".."},
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		for _, test := range relativePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			expected := path.Join(cwd, test.out)

			if got != expected {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, expected, got)
			}
		}
	})

	t.Run("WithTildePath", func(t *testing.T) {
		var tildePathInputs = []struct {
			in  string
			out string
		}{
			{"~", ""},
			{"~/", ""},
			{"~/xyz", "xyz"},
			{"~/xyz/./abc", "xyz/abc"},
			{"~/xyz/../abc", "abc"},
			{"~/xyz/abc/..", "xyz"},
			{"~/xyz/../abc/..", ""},
			{"~/xyz/../..", ".."},
			{"~/xyz/~/abc", "xyz/~/abc"},
		}

		usr, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		dir := usr.HomeDir

		for _, test := range tildePathInputs {
			got, _ := filesystem.AbsolutePath(test.in)
			expected := path.Join(dir, test.out)

			if got != expected {
				t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, expected, got)
			}
		}
	})
}

func TestIsFile(t *testing.T) {
	var objectsTest = []struct {
		in  string
		out bool
	}{
		{"testdata/file.txt", true},
		{"testdata/folder", false},
		{"testdata/non-existent-file", false},
	}

	for _, test := range objectsTest {
		got := filesystem.IsFile(test.in)
		if got != test.out {
			t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, test.out, got)
		}
	}
}

func TestIsDir(t *testing.T) {
	var objectsTest = []struct {
		in  string
		out bool
	}{
		{"testdata/file.txt", false},
		{"testdata/folder", true},
		{"testdata/non-existent-dir", false},
	}

	for _, test := range objectsTest {
		got := filesystem.IsDir(test.in)
		if got != test.out {
			t.Errorf("failed for '%s': expected '%v', got '%v'", test.in, test.out, got)
		}
	}
}

func TestRelativePath(t *testing.T) {
	var objectsTest = []struct {
		base string
		in   string
		out  string
	}{
		{base: "/foo/bar", in: "/foo/bar/xyz", out: "xyz"},
		{base: "/foo/bar/", in: "/foo/bar/xyz", out: "xyz"},
		{base: "/foo/bar", in: "/foo/bar/xyz/", out: "xyz"},
		{base: "/foo/bar", in: "foo/bar/xyz", out: "foo/bar/xyz"},
		{base: "/foo/bar", in: "/foo/bar", out: "."},
		{base: "/foo/bar/", in: "/foo/bar", out: "."},
		{base: "/foo/bar", in: "/foo/bar/", out: "."},
		{base: "", in: "/foo/bar", out: "/foo/bar"},
		{base: "/foo/bar", in: "/abc/def", out: "/abc/def"},
	}
	for _, tt := range objectsTest {
		got := filesystem.RelativePath(tt.base, tt.in)
		if got != tt.out {
			t.Errorf("failed for base '%s', path '%s': expected '%s', got '%s'", tt.base, tt.in, tt.out, got)
		}
	}
}

func TestEmptyDirWithOneFile(t *testing.T) {
	defer filet.CleanUp(t)

	// create a dir with a file inside of it
	dir := filet.TmpDir(t, "")
	file := filet.TmpFile(t, dir, "")

	err := filesystem.EmptyDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	// dir should exists after removal
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("failed: dir has been deleted, that was unexpected")
	}

	// file should not exists after removal
	if _, err := os.Stat(file.Name()); err == nil {
		t.Errorf("failed: there are content inside dir, that was unexpected")
	}
}

func TestEmptyDirWithOneDir(t *testing.T) {
	defer filet.CleanUp(t)

	// create a dir with a file inside of it
	parentDir := filet.TmpDir(t, "")
	childDir := filet.TmpDir(t, parentDir)

	err := filesystem.EmptyDir(parentDir)
	if err != nil {
		t.Fatal(err)
	}

	// parent dir should exists after removal
	if _, err := os.Stat(parentDir); err != nil {
		t.Errorf("failed: paret dir has been deleted, that was unexpected")
	}

	// child dir should not exists after removal
	if _, err := os.Stat(childDir); err == nil {
		t.Errorf("failed: there are content inside parent dir, that was unexpected")
	}
}

func TestEmptyOrCreateDirExistingDir(t *testing.T) {
	const numberOfIterations = 10
	defer filet.CleanUp(t)

	for i := 0; i < numberOfIterations; i++ {
		// create a dir with a file inside of it
		dir := filet.TmpDir(t, "")
		file := filet.TmpFile(t, dir, "")

		err := filesystem.EmptyOrCreateDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		// dir should exists after removal
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("failed: dir has been deleted, that was unexpected")
		}

		// file should not exists after removal
		if _, err := os.Stat(file.Name()); err == nil {
			t.Errorf("failed: there are content inside dir, that was unexpected")
		}
	}
}

func TestEmptyOrCreateDirNonExistingDir(t *testing.T) {
	const numberOfIterations = 10
	var dirs [numberOfIterations]string

	for i := 0; i < numberOfIterations; i++ {
		dirs[i] = filepath.Join(os.TempDir(), fmt.Sprintf("dir-%d.%d", i, time.Now().UnixNano()))
	}
	defer func(dirs [numberOfIterations]string) {
		for _, dir := range dirs {
			_ = os.RemoveAll(dir)
		}
	}(dirs)

	for _, dir := range dirs {
		err := filesystem.EmptyOrCreateDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		// dir should exists
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("failed: dir has been deleted, that was unexpected")
		}
	}
}
