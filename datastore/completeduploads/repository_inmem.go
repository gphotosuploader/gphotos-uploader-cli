package completeduploads

type InMemRepository struct {
	m map[string]*CompletedUploadedFileItem
}

// NewInMemRepository create a new repository
func NewInMemRepository() *InMemRepository {
	var m = map[string]*CompletedUploadedFileItem{}
	return &InMemRepository{m: m}
}

// Get an item
func (r *InMemRepository) Get(path string) (CompletedUploadedFileItem, error) {
	if r.m[path] == nil {
		return CompletedUploadedFileItem{}, ErrNotFound
	}
	return *r.m[path], nil
}

// Store an item
func (r *InMemRepository) Put(item CompletedUploadedFileItem) error {
	r.m[item.path] = &item
	return nil
}

// Delete an item
func (r *InMemRepository) Delete(path string) error {
	if r.m[path] == nil {
		return ErrNotFound
	}
	r.m[path] = nil
	return nil
}
