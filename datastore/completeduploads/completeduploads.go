package completeduploads

import (
	"fmt"

	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

// Service represents the repository where uploaded objects are tracked
type Service struct {
	repo Repository
}

// NewService created a Service to track uploaded objects
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// Close closes the service.
//
// No operation could be done after that.
func (s *Service) Close() error {
	return s.repo.Close()
}

// IsAlreadyUploaded checks if the file was already uploaded
func (s *Service) IsAlreadyUploaded(filePath string) (bool, error) {
	// find a previous upload in the repository
	item, err := s.repo.Get(filePath)
	if err != nil {
		// this file was not uploaded before
		return false, nil
	}

	// value found on the cache

	// get the last modified time from the cache
	cacheMtime, err := item.GetTrackedMTime()
	if err != nil {
		return false, err
	}

	// check stored last modified time with the current one to see if the
	// file has been modified
	if cacheMtime != 0 {
		fileMtime, err := filesystem.GetMTime(filePath)
		if err != nil {
			return false, err
		}
		if fileMtime.Unix() == cacheMtime {
			return true, nil
		}
	}

	// file was not uploaded before or modified time has changed after being
	// uploaded
	fileHash, err := Hash(filePath)
	if err != nil {
		return false, err
	}

	// checks if the file is the same (equal hash)
	if item.GetTrackedHash() == fmt.Sprint(fileHash) {
		// update last modified time on the cache
		err = s.CacheAsAlreadyUploaded(filePath)
		if err != nil {
			return true, err
		}

	}

	return false, nil
}

// CacheAsAlreadyUploaded marks a file as already uploaded to prevent re-uploads
func (s *Service) CacheAsAlreadyUploaded(filePath string) error {
	item, err := NewCompletedUploadedFileItem(filePath)
	if err != nil {
		return err
	}
	return s.repo.Put(item)
}

// RemoveAsAlreadyUploaded removes a file previously marked as uploaded
func (s *Service) RemoveAsAlreadyUploaded(filePath string) error {
	return s.repo.Delete(filePath)
}
