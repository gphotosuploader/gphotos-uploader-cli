package filetracker

import "github.com/syndtr/goleveldb/leveldb"

type LevelDBRepository struct {
	db *leveldb.DB
}

// NewLevelDBRepository creates a repository using LevelDB.
func NewLevelDBRepository(db *leveldb.DB) *LevelDBRepository {
	return &LevelDBRepository{db: db}
}

// Get returns the item specified by key. It returns ErrItemNotFound if the
// DB does not contains the key.
func (r *LevelDBRepository) Get(key string) (TrackedFile, error) {
	val, err := r.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return TrackedFile{}, ErrItemNotFound
	}
	if err != nil {
		return TrackedFile{}, err
	}
	return NewTrackedFile(string(val)), nil
}

// Put stores the item under key.
func (r *LevelDBRepository) Put(key string, item TrackedFile) error {
	return r.db.Put([]byte(key), []byte(item.value), nil)
}

// Delete removes the item specified by key.
func (r *LevelDBRepository) Delete(key string) error {
	if err := r.db.Delete([]byte(key), nil); err != nil {
		return err
	}
	return nil
}

// Close closes the DB.
func (r *LevelDBRepository) Close() error {
	return r.db.Close()
}
