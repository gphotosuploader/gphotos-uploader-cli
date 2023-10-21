package filetracker

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
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

// Close closes the DB.
func (r LevelDBRepository) Close() error {
	return r.DB.Close()
}

// Destroy completely remove an existing LevelDB database directory.
func (r LevelDBRepository) Destroy() error {
	_ = r.DB.Close()
	return os.RemoveAll(r.path)
}
