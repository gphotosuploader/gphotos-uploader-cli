package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/palantir/stacktrace"
)

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

// GetLastModifiedTime returns the Last Modified time from the file
func GetMTime(path string) (mtime time.Time, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	return
}

// BufferFromFile opens the file to return a buffer
func BufferFromFile(filePath string) (buf []byte, _ error) {
	if !IsFile(filePath) {
		return nil, fmt.Errorf("not a file")
	}
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", filePath)
	}

	return buf, nil
}
