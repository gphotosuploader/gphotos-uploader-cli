package filetracker

import (
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// FileTracker allows tracking already uploaded files in a repository.
type FileTracker struct {
	repo FileRepository

	// Hasher allows changing the way that hashes are calculated.
	// Uses xxHash32Hasher{} by default.
	// Useful for testing.
	Hasher Hasher

	Logger log.Logger
}

// Hasher is a Hasher to get the value of the file.
type Hasher interface {
	Hash(file string) (string, error)
}

// FileRepository is the repository where to track already uploaded files.
type FileRepository interface {
	Get(key string) (TrackedFile, bool)
	Put(key string, item TrackedFile) error
	Delete(key string) error
	Close() error
	Destroy() error
}

// New returns a FileTracker using specified repo.
func New(r FileRepository) *FileTracker {
	return &FileTracker{
		repo:   r,
		Hasher: XXHash32Hasher{},
		Logger: log.Discard,
	}
}

// MarkAsUploaded marks a file as already uploaded.
func (ft FileTracker) MarkAsUploaded(file string) error {
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
		Hash:    hash,
	}

	return ft.repo.Put(file, item)
}

// IsUploaded checks if the file was already uploaded.
// First compares the last modification time of the file against the one in the repository.
// Last time modification comparison tries to reduce the number of times when the hash comparison
// is needed.
// In case that last modification time has changed (or it doesn't exist - retro compatibility),
// it compares a hash of the content of the file against the one in the repository.
func (ft FileTracker) IsUploaded(file string) bool {
	item, found := ft.repo.Get(file)
	if !found {
		return false
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		ft.Logger.Debugf("Error retrieving file info for '%s' (%s).", file, err)
		return false
	}

	if item.ModTime.Equal(fileInfo.ModTime()) {
		return true
	}

	hash, err := ft.Hasher.Hash(file)
	if err != nil {
		return false
	}

	// checks if the file is the same (equal value)
	if item.Hash == hash {
		// updates file marker with mtime to speed up comparison on the next run
		item.ModTime = fileInfo.ModTime()
		if err = ft.repo.Put(file, item); err != nil {
			ft.Logger.Debugf("Error updating marker for '%s' with modification time (%s).", file, err)
		}

		return true
	}

	return false
}

// UnmarkAsUploaded un-marks a file as already uploaded.
func (ft FileTracker) UnmarkAsUploaded(file string) error {
	return ft.repo.Delete(file)
}

// Close closes the file tracker repository.
// No operation could be done after that.
func (ft FileTracker) Close() error {
	return ft.repo.Close()
}

// Destroy completely remove an existing FileTracker database.
func (ft FileTracker) Destroy() error {
	return ft.repo.Destroy()
}
