// Package upload_tracker provides implementation of [gphotosuploader/google-photos-api-client-go] Store interface
package upload_tracker

import (
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

type LevelDBStore struct {
	db   *leveldb.DB
	path string
}

// NewStore create a new Store implemented by LevelDB
func NewStore(path string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	s := &LevelDBStore{
		db:   db,
		path: path,
	}
	return s, err
}

// Get returns the value corresponding to the given key
func (s *LevelDBStore) Get(key string) (string, bool) {
	v, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return "", false
	}
	return string(v), true
}

// Set stores the url for a given fingerprint
func (s *LevelDBStore) Set(key string, value string) {
	_ = s.db.Put([]byte(key), []byte(value), nil)
}

func (s *LevelDBStore) Delete(key string) {
	_ = s.db.Delete([]byte(key), nil)
}

// Close closes the service
func (s *LevelDBStore) Close() {
	_ = s.db.Close()
}

// Destroy completely remove an existing LevelDB database directory.
func (s *LevelDBStore) Destroy() error {
	_ = s.db.Close()
	return os.RemoveAll(s.path)
}
