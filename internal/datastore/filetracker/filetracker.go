package filetracker

import (
	"fmt"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

var (
	// ErrItemNotFound is the expected error if the item is not found.
	ErrItemNotFound = fmt.Errorf("item was not found")
)

// FileTracker allows to track already uploaded files in a repository.
type FileTracker struct {
	repo Repository

	// Hasher allows to change the way that hashes are calculated. Uses xxHash32Hasher{} by default.
	// Useful for testing.
	Hasher Hasher

	logger log.Logger
}

// Hasher is a Hasher to get the value of the file.
type Hasher interface {
	Hash(file string) (string, error)
}

// Repository is the repository where to track already uploaded files.
type Repository interface {
	// Get It returns ErrItemNotFound if the repo does not contains the key.
	Get(key string) (TrackedFile, error)
	Put(key string, item TrackedFile) error
	Delete(key string) error
	Close() error
}

// New returns a FileTracker using specified repo.
func New(r Repository) *FileTracker {
	return &FileTracker{
		repo:   r,
		Hasher: xxHash32Hasher{},
		logger: log.GetInstance(),
	}
}

// Put marks a file as already uploaded to prevent re-uploads.
func (ft FileTracker) Put(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	hash, err := ft.Hasher.Hash(file)
	if err != nil {
		return err
	}
	item := TrackedFile{
		ModTime: fileInfo.ModTime(),
		Hash: hash,
	}

	return ft.repo.Put(file, item)
}

// Exist checks if the file was already uploaded.
// Exist compares the value of the file against the repository.
func (ft FileTracker) Exist(file string) bool {
	// Get returns ErrItemNotFound if the repo does not contains the key.
	item, err := ft.repo.Get(file)
	if err != nil {
		return false
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		ft.logger.Debugf("Error retrieving file info for '%s' (%s).", file, err)
		return false
	}

	if item.ModTime.Equal(fileInfo.ModTime()) {
		ft.logger.Debugf("File modification time has not changed for '%s'.", file)
		return true
	}

	hash, err := ft.Hasher.Hash(file)
	if err != nil {
		return false
	}

	// checks if the file is the same (equal value)
	if item.Hash == hash {
		ft.logger.Debugf("File hash has not changed for '%s'.", file)
		return true
	}

	return false
}

// Delete un-marks a file as already uploaded.
func (ft FileTracker) Delete(file string) error {
	return ft.repo.Delete(file)
}

// Close closes the file tracker repository.
// No operation could be done after that.
func (ft FileTracker) Close() error {
	return ft.repo.Close()
}
