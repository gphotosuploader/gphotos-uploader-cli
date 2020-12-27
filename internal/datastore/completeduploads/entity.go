package completeduploads

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pierrec/xxHash/xxHash32"
)

var (
	// ErrNotFound not found
	ErrNotFound = fmt.Errorf("not found")

	// ErrCannotBeDeleted bookmark cannot be deleted
	ErrCannotBeDeleted = fmt.Errorf("cannot be deleted")

	// ErrCannotGetMTime
	ErrCannotGetMTime = fmt.Errorf("failed getting local image mtime")
)

type CompletedUploadedFileItem struct {
	path  string
	value string
}

// NewCompletedUploadedFileItem creates a new item for the specified file
func NewCompletedUploadedFileItem(filePath string) (CompletedUploadedFileItem, error) {
	item := CompletedUploadedFileItem{
		path: filePath,
	}

	fileHash, err := Hash(filePath)
	if err != nil {
		return item, err
	}

	mTime, err := getMTime(filePath)
	if err != nil {
		return item, ErrCannotGetMTime
	}

	item.SetValue(fileHash, mTime)
	return item, nil
}

// SetValue stores the Hash and mTime in the k/v store.
func (f *CompletedUploadedFileItem) SetValue(hash uint32, mTime time.Time) {
	f.value = strconv.FormatInt(mTime.Unix(), 10) + "|" + fmt.Sprint(hash)
}

// GetTrackedHash returns the hash value stored in the cache.
func (f *CompletedUploadedFileItem) GetTrackedHash() string {
	parts := strings.Split(f.value[:], "|")
	if len(parts) > 1 {
		return parts[1]
	}
	return parts[0]
}

// GetTrackedMTime returns the last modified time value stored in the cache.
func (f *CompletedUploadedFileItem) GetTrackedMTime() (int64, error) {
	parts := strings.Split(f.value[:], "|")
	if len(parts) <= 1 {
		return 0, ErrCannotGetMTime
	}

	cacheMtime, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return cacheMtime, nil
}

// Hash return the hash of a file
func Hash(filePath string) (uint32, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	hasher := xxHash32.New(0xCAFE) // value.Hash32
	defer hasher.Reset()

	_, err = io.Copy(hasher, inputFile)
	if err != nil {
		return 0, err
	}

	return hasher.Sum32(), nil
}
