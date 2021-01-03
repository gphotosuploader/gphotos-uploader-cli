package mock

// FileTracker represents a service to track already uploaded files.
type FileTracker struct {
	CacheAsAlreadyUploadedFn      func(path string) error
	CacheAsAlreadyUploadedInvoked bool

	IsAlreadyUploadedFn      func(path string) (bool, error)
	IsAlreadyUploadedInvoked bool

	RemoveAsAlreadyUploadedFn      func(path string) error
	RemoveAsAlreadyUploadedInvoked bool
}

func (t *FileTracker) CacheAsAlreadyUploaded(path string) error {
	t.CacheAsAlreadyUploadedInvoked = true
	return t.CacheAsAlreadyUploadedFn(path)
}

func (t *FileTracker) IsAlreadyUploaded(path string) (bool, error) {
	t.IsAlreadyUploadedInvoked = true
	return t.IsAlreadyUploadedFn(path)
}

func (t *FileTracker) RemoveAsAlreadyUploaded(path string) error {
	t.RemoveAsAlreadyUploadedInvoked = true
	return t.RemoveAsAlreadyUploadedFn(path)
}
