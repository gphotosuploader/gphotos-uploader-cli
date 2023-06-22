package mock

// FileTracker mocks the service to track already uploaded files.
type FileTracker struct {
	MarkAsUploadedFn   func(path string) error
	IsUploadedFn       func(path string) bool
	UnmarkAsUploadedFn func(path string) error
}

// MarkAsUploaded invokes the mock implementation.
func (t *FileTracker) MarkAsUploaded(path string) error {
	return t.MarkAsUploadedFn(path)
}

// IsUploaded invokes the mock implementation.
func (t *FileTracker) IsUploaded(path string) bool {
	return t.IsUploadedFn(path)
}

// UnMarkAsUploaded invokes the mock implementation.
func (t *FileTracker) UnmarkAsUploaded(path string) error {
	return t.UnmarkAsUploadedFn(path)
}
