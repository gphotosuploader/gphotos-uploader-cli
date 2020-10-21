package upload

// UploadFolderJob represents a job to upload all photos from the specified folder
type UploadFolderJob struct {
	FileTracker FileTracker

	SourceFolder       string
	CreateAlbum        bool
	CreateAlbumBasedOn string
	Filter             FileFilterer
}

// FileTracker represents a service to track already uploaded files.
type FileTracker interface {
	CacheAsAlreadyUploaded(filePath string) error
	IsAlreadyUploaded(filePath string) (bool, error)
	RemoveAsAlreadyUploaded(filePath string) error
}

// FileFilterer represents a way to implement include/exclude files filtering.
type FileFilterer interface {
	IsAllowed(path string) bool
	IsExcluded(path string) bool
}
