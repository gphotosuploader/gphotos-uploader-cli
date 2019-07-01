package completeduploads

// Repository represents a database where to track uploaded files
type Repository interface {
	Get(key string) (CompletedUploadedFileItem, error)
	Put(item CompletedUploadedFileItem) error
	Delete(key string) error
}
