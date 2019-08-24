package uploadurls

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pierrec/xxHash/xxHash32"
	"github.com/syndtr/goleveldb/leveldb"
)

type Service struct {
	db *leveldb.DB
}

func NewService(db *leveldb.DB) *Service {
	return &Service{db}
}

// TODO: refactor with completeduploads
func hashFile(filePath string) (uint32, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	hasher := xxHash32.New(0xCAFE) // hash.Hash32
	defer hasher.Reset()

	_, err = io.Copy(hasher, inputFile)
	if err != nil {
		return 0, err
	}

	return hasher.Sum32(), nil
}

// GetUploadURL gets upload URL from database if available for resumable upload
func (s *Service) GetUploadURL(filePath string) (string, error) {
	// look for upload URL in database
	val, err := s.db.Get([]byte(filePath), nil)
	if err == leveldb.ErrNotFound {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	// value found, try to split hash and upload URL
	strval := string(val[:])
	parts := strings.Split(strval, "|")
	if len(parts) != 2 {
		return "", fmt.Errorf("failed parsing upload URL data from database: %s", strval)

	}
	cacheHash := parts[0]
	uploadURL := parts[1]

	fileHash, err := hashFile(filePath)
	if err != nil {
		return "", err
	}

	if cacheHash != fmt.Sprint(fileHash) {
		err = s.RemoveUploadURL(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to remove upload URL from database: %s", err)
		}

		return "", nil
	}

	return uploadURL, err
}

// PutUploadURL puts a file's upload URL in database for resumable upload
func (s *Service) PutUploadURL(filePath, uploadURL string) error {
	fileHash, err := hashFile(filePath)
	if err != nil {
		return err
	}

	val := fmt.Sprint(fileHash) + "|" + uploadURL
	err = s.db.Put([]byte(filePath), []byte(val), nil)
	if err != nil {
		return err
	}
	log.Printf("Stored upload URL for file %s (%s)", filePath, uploadURL)

	return nil
}

// RemoveUploadURL removes a file's upload URL from the database
func (s *Service) RemoveUploadURL(filePath string) error {
	log.Printf("Removing file's upload URL from DB: %s", filePath)
	return s.db.Delete([]byte(filePath), nil)
}
