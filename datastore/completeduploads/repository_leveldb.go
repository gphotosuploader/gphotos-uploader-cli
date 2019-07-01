package completeduploads

import "github.com/syndtr/goleveldb/leveldb"

type LevelDBRepository struct {
	db *leveldb.DB
}

// NewLevelDBRepository create a new repository
func NewLevelDBRepository(db *leveldb.DB) *LevelDBRepository {
	return &LevelDBRepository{db: db}
}

// Get an item
func (r *LevelDBRepository) Get(path string) (CompletedUploadedFileItem, error) {
	val, err := r.db.Get([]byte(path), nil)
	if err != nil {
		return CompletedUploadedFileItem{}, err
	}

	item := CompletedUploadedFileItem{
		path:  path,
		value: string(val),
	}
	return item, nil
}

// Store an item
func (r *LevelDBRepository) Put(item CompletedUploadedFileItem) error {
	return r.db.Put([]byte(item.path), []byte(item.value), nil)
}

// Delete an item
func (r *LevelDBRepository) Delete(path string) error {
	err := r.db.Delete([]byte(path), nil)
	if err != nil {
		return ErrCannotBeDeleted
	}
	return nil
}
