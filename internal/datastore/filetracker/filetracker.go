package filetracker

import (
	"fmt"
	"os"
	"time"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// FileTracker allows tracking already uploaded files in a repository.
type FileTracker struct {
	repo FileRepository

	// Hasher allows changing the way that hashes are calculated.
	// Uses SHA256Hasher{} by default.
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
	All() (map[string]TrackedFile, error)
	Close() error
	Destroy() error
}

func (ft FileTracker) Receipt(file string) (TrackedFile, bool)   { return ft.repo.Get(file) }
func (ft FileTracker) Receipts() (map[string]TrackedFile, error) { return ft.repo.All() }

// New returns a FileTracker using specified repo.
func New(r FileRepository) *FileTracker {
	return &FileTracker{
		repo:   r,
		Hasher: SHA256Hasher{},
		Logger: log.Discard,
	}
}

// UploadReceipt is the durable proof returned after Google creates a media item.
type UploadReceipt struct {
	MediaItemID string
	ProductURL  string
	SHA256      string
	Size        int64
	UploadedAt  time.Time
}

// RecordUpload atomically associates verified local bytes with Google's media item ID.
func (ft FileTracker) RecordUpload(file string, receipt UploadReceipt) error {
	if receipt.MediaItemID == "" || receipt.SHA256 == "" || receipt.Size < 0 {
		return fmt.Errorf("incomplete upload receipt")
	}
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	if info.Size() != receipt.Size {
		return fmt.Errorf("file size changed after upload")
	}
	hash, err := ft.Hasher.Hash(file)
	if err != nil {
		return err
	}
	if hash != receipt.SHA256 {
		return fmt.Errorf("file content changed after upload")
	}
	if receipt.UploadedAt.IsZero() {
		receipt.UploadedAt = time.Now().UTC()
	}
	return ft.repo.Put(file, TrackedFile{Version: 2, ModTime: info.ModTime(), Size: info.Size(), Hash: hash, MediaItemID: receipt.MediaItemID, ProductURL: receipt.ProductURL, UploadedAt: receipt.UploadedAt})
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
	if item.Version >= 2 {
		if item.Size != fileInfo.Size() {
			return false
		}
		hash, err := ft.Hasher.Hash(file)
		return err == nil && item.Hash == hash
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
