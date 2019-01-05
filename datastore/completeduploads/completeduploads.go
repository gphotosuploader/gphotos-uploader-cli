package completeduploads

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"github.com/pierrec/xxHash/xxHash32"
	"github.com/syndtr/goleveldb/leveldb"
)

type CompletedUploadsService struct {
	db *leveldb.DB
}

func NewService(db *leveldb.DB) *CompletedUploadsService {
	return &CompletedUploadsService{db}
}

func fileHash(filePath string) (uint32, error) {
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

//func uint32ToBytes(u uint32) []byte {
//	a := make([]byte, 4)
//	binary.LittleEndian.PutUint32(a, u)
//	return a
//}

// IsAlreadyUploaded checks in cache if the file was already uploaded
func (s *CompletedUploadsService) IsAlreadyUploaded(filePath string) (bool, error) {
	isUploaded := false

	// look for previous upload in cache
	val, err := s.db.Get([]byte(filePath), nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}

	if err == nil {
		// value found, try to split mtime and hash
		parts := strings.Split(string(val[:]), "|")
		cacheMtime := int64(0)
		cacheHash := ""
		if len(parts) > 1 {
			cacheMtime, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return false, err
			}
			cacheHash = parts[1]
		} else {
			cacheHash = parts[0]
		}
		// check mtime first
		if cacheMtime != 0 {
			fileMtime, err := filesystem.GetMTime(filePath)
			if err != nil {
				return false, err
			}
			if fileMtime.Unix() == cacheMtime {
				isUploaded = true
				//log.Printf("%s mtime matched %i", filePath, cacheMtime)
			}
		}
		// mtime is different, check hash
		if !isUploaded {
			fileHash, err := fileHash(filePath)
			if err != nil {
				return false, err
			}

			if cacheHash == fmt.Sprint(fileHash) {
				isUploaded = true
				//log.Printf("%s hash match %s", filePath, cacheHash)
				// update db mtime
				err = s.CacheAsAlreadyUploaded(filePath)
				if err != nil {
					return isUploaded, err
				}
			}
		}
	}

	return isUploaded, err
}

// CacheAsAlreadyUploaded marks a file in cache as already uploaded to prevent re-uploads
func (s *CompletedUploadsService) CacheAsAlreadyUploaded(filePath string) error {
	fileHash, err := fileHash(filePath)
	if err != nil {
		return err
	}

	mtime, err := filesystem.GetMTime(filePath)
	if err != nil {
		return fmt.Errorf("failed getting local image mtime")
	}

	val := strconv.FormatInt(mtime.Unix(), 10) + "|" + fmt.Sprint(fileHash)
	err = s.db.Put([]byte(filePath), []byte(val), nil)
	if err != nil {
		return err
	}
	log.Printf("Marked as uploaded: %s", filePath)

	return nil
}

// RemoveAsAlreadyUploaded removes a file previously marked as uploaded from the db
func (s *CompletedUploadsService) RemoveAsAlreadyUploaded(filePath string) error {
	log.Printf("Removing file from upload DB: %s", filePath)
	err := s.db.Delete([]byte(filePath), nil)

	return err
}
