package filesystem

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// AbsolutePath converts a path (relative or absolute) into an absolute one.
// Supports '~' notation for $HOME directory of the current user.
func AbsolutePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		return dir, nil
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		return filepath.Join(dir, path[2:]), nil
	}
	return filepath.Abs(path)
}

// EmptyDir removes all files and folders inside the specified path.
// It could be similar to RemoveAll() but without removing the folder itself.
func EmptyDir(path string) error {
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDirIfDoesNotExist creates a directory if the specified path does not exist
func CreateDirIfDoesNotExist(path string) error {
	if IsDir(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// EmptyOrCreateDir create a new folder or empty an existing one
func EmptyOrCreateDir(path string) error {
	if err := CreateDirIfDoesNotExist(path); err != nil {
		return err
	}
	return EmptyDir(path)
}

// RelativePath returns a path relative to the given base path. If the path is not
// under the given base path, the specified path is returned. So all returned paths
// are under the base path.
func RelativePath(basepath, path string) string {
	rp, err := filepath.Rel(basepath, path)
	if err != nil {
		return path
	}
	if strings.HasPrefix(rp, "../") {
		return path
	}
	return rp
}

// IsFile asserts there is a file at path
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

// IsDir asserts there is a directory at path
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

// GetMTime returns the Last Modified time from the file
func GetMTime(path string) (mtime time.Time, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	return
}
