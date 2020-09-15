package upload


// UploadFolderJob represents a job to upload all photos from the specified folder
type UploadFolderJob struct {
	FileTracker   FileTracker

	SourceFolder       string
	CreateAlbum        bool
	CreateAlbumBasedOn string
	Filter             *Filter
}
