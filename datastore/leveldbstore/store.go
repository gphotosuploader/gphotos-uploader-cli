// Package leveldbstore provides implementation of LevelDB key/value database.
//
// Create or open a database:
//
//	// The returned DB instance is safe for concurrent use. Which mean that all
//	// DB's methods may be called concurrently from multiple goroutine.
//	db, err := leveldbstore.NewStore("path/to/db")
//	...
//	defer db.Close()
//	...
//
// Read or modify the database content:
//
//	// Remember that the contents of the returned slice should not be modified.
//	data := db.Get(key)
//	...
//	db.Put(key), []byte("value"))
//	...
//	db.Delete(key)
//	...
package leveldbstore

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBStore struct {
	db *leveldb.DB
}

// NewStore create a new Store implemented by LevelDB
func NewStore(path string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	s := &LevelDBStore{db: db}
	return s, err
}

// Get returns the value corresponding to the given key
func (s *LevelDBStore) Get(key string) []byte {
	v, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return []byte{}
	}
	return v
}

// Set stores the url for a given fingerprint
func (s *LevelDBStore) Set(key string, value []byte) {
	_ = s.db.Put([]byte(key), value, nil)
}

func (s *LevelDBStore) Delete(key string) {
	_ = s.db.Delete([]byte(key), nil)
}

func (s *LevelDBStore) Close() {
	s.Close()
}
