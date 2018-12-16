package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/palantir/stacktrace"
)

const (
	HEADERSIZE = 100
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
		return nil, stacktrace.Propagate(err, "Failed reading file: %s: Ignoring file...\n", filePath)
	}

	return buf, nil
}

// BufferHeaderFromFile opens the file to return a buffer of the first HEADERSIZE bytes
func BufferHeaderFromFile(filePath string) ([]byte, error) {
	if !IsFile(filePath) {
		return nil, fmt.Errorf("not a file")
	}
	r, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	var buf [HEADERSIZE]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed reading %s bytes of file: %s: Ignoring file...\n", strconv.FormatInt(HEADERSIZE, 10), filePath)
	}

	return buf, nil
}
