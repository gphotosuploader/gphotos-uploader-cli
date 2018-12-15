package completeduploads

import (
	"encoding/binary"
	"fmt"

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
	fileBuf, err := filesystem.BufferFromFile(filePath)
	if err != nil {
		return 0, err
	}

	hasher := xxHash32.New(0xCAFE) // hash.Hash32
	defer hasher.Reset()

	hasher.Write(fileBuf)
	fmt.Printf("%x\n", hasher.Sum32())

	return hasher.Sum32(), nil
}

func uint32ToBytes(u uint32) []byte {
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, u)
	return a
}

// IsAlreadyUploaded checks in cache if the file was already uploaded
func (s *CompletedUploadsService) IsAlreadyUploaded(filePath string) (bool, error) {
	fileHash, err := fileHash(filePath)
	if err != nil {
		return false, err
	}

	cacheVal, err := s.db.Get(uint32ToBytes(fileHash), nil)
	if err != nil {
		return false, err
	}

	if cacheVal != nil {
		return true, nil
	}

	return false, nil
}

// CacheAsAlreadyUploaded marks a file in cache as already uploaded to prevent re-uploads
func (s *CompletedUploadsService) CacheAsAlreadyUploaded(filePath string) error {
	fileHash, err := fileHash(filePath)
	if err != nil {
		return err
	}

	err = s.db.Put(uint32ToBytes(fileHash), []byte{}, nil)
	if err != nil {
		return err
	}

	return nil
}

// func IsUploadedPrev(filePath string, db *leveldb.DB) (bool, error) {
// 	isUploaded := false

// 	// look for previous upload in cache
// 	val, err := db.Get([]byte(filePath), nil)
// 	if err == nil {
// 		// value found, try to split mtime and hash
// 		parts := strings.Split(string(val[:]), "|")
// 		cacheMtime := int64(0)
// 		cacheHash := ""
// 		if len(parts) > 1 {
// 			cacheMtime, err = strconv.ParseInt(parts[0], 10, 64)
// 			cacheHash = parts[1]
// 		} else {
// 			cacheHash = parts[0]
// 		}
// 		// check mtime first
// 		if err == nil && cacheMtime != 0 {
// 			fileMtime, err := filesystem.GetMTime(filePath)
// 			if err == nil && fileMtime.Unix() == cacheMtime {
// 				isUploaded = true
// 				//log.Printf("%s mtime matched %i", filePath, cacheMtime)
// 			}
// 		}
// 		// mtime is different, check hash
// 		if !isUploaded {
// 			localImg, err := imageFromPath(filePath)
// 			if err != nil {
// 				err = fmt.Errorf("failed loading local image from path")
// 			} else {
// 				localHash := getImageHash(localImg)
// 				isUploaded = isSameHash(cacheHash, localHash)
// 				if isUploaded {
// 					//log.Printf("%s hash match %s", filePath, cacheHash)
// 					// update db mtime
// 					err = MarkUploaded(filePath, db)
// 				}
// 			}
// 		}
// 	} else if strings.Contains(err.Error(), "not found") {
// 		err = nil
// 	}

// 	return isUploaded, err
// }

// func MarkUploaded(filePath string, db *leveldb.DB) error {
// 	localImg, err := imageFromPath(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed loading local image from path")
// 	}
// 	mtime, err := filesystem.GetMTime(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed getting local image mtime")
// 	}
// 	val := strconv.FormatInt(mtime.Unix(), 10) + "|" + getImageHash(localImg)

// 	log.Printf("Marking file as uploaded: %s", filePath, val)
// 	err = db.Put([]byte(filePath), []byte(val), nil)

// 	return err
// }

// func removeFromDB(filePath string, db *leveldb.DB) error {
// 	log.Printf("Removing file from upload DB: %s", filePath)
// 	err := db.Delete([]byte(filePath), nil)

// 	return err
// }
