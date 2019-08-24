package filesystem

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// AbsolutePath converts a path (relative or absolute) into an absolute one.
// Supports '~' notation for $HOME directory of the current user.
func AbsolutePath(path string) string {
	usr, err := user.Current()
	if err != nil {
		// It's very strange to have an error here, but in any case,
		// return the same path was given.
		return path
	}
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		return dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		return filepath.Join(dir, path[2:])
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		log.Println(err)
	}
	return abspath
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
