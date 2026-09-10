package filetracker

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"os"
)

// DB represents a LevelDB database.
type DB interface {
	Get(key []byte, ro *opt.ReadOptions) ([]byte, error)
	Put(key []byte, item []byte, wo *opt.WriteOptions) error
	Delete(key []byte, wo *opt.WriteOptions) error
	Close() error
}

// LevelDBRepository implements a FileRepository using LevelDB.
type LevelDBRepository struct {
	DB   DB
	path string
}

// NewLevelDBRepository creates a repository using LevelDB package.
func NewLevelDBRepository(path string) (*LevelDBRepository, error) {
	ft, err := leveldb.OpenFile(path, nil)
	return &LevelDBRepository{
		DB:   ft,
		path: path,
	}, err
}

// Get returns the item specified by key. It returns ErrItemNotFound if the
// DB does not contain the key.
func (r LevelDBRepository) Get(key string) (TrackedFile, bool) {
	val, err := r.DB.Get([]byte(key), nil)
	if err != nil {
		return TrackedFile{}, false
	}
	return NewTrackedFile(string(val)), true
}

// Put stores the item under key.
func (r LevelDBRepository) Put(key string, item TrackedFile) error {
	return r.DB.Put([]byte(key), []byte(item.String()), nil)
}

// Delete removes the item specified by key.
func (r LevelDBRepository) Delete(key string) error {
	return r.DB.Delete([]byte(key), nil)
}

// All returns a consistent snapshot of every tracked upload receipt.
func (r LevelDBRepository) All() (map[string]TrackedFile, error) {
	result := make(map[string]TrackedFile)
	db, ok := r.DB.(interface {
		NewIterator(*util.Range, *opt.ReadOptions) iterator.Iterator
	})
	if !ok {
		return nil, fmt.Errorf("repository does not support iteration")
	}
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		result[string(iter.Key())] = NewTrackedFile(string(iter.Value()))
	}
	return result, iter.Error()
}

// Close closes the DB.
func (r LevelDBRepository) Close() error {
	return r.DB.Close()
}

// Destroy completely remove an existing LevelDB database directory.
func (r LevelDBRepository) Destroy() error {
	_ = r.DB.Close()
	return os.RemoveAll(r.path)
}
