package config

// Config represents the application settings.
type Config struct {
	ConfigPath         string
	SecretsBackendType string
	APIAppCredentials  *APIAppCredentials
	Jobs               []FolderUploadJob
}

// APIAppCredentials represents Google Photos API credentials for OAuth
type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

// FolderUploadJob represents configuration for a folder to be uploaded
type FolderUploadJob struct {
	Account           string
	SourceFolder      string
	MakeAlbums        MakeAlbums
	DeleteAfterUpload bool
	UploadVideos      bool
	IncludePatterns   []string
	ExcludePatterns   []string
}

// MakeAlbums represents configuration about how to create Albums in Google Photos
type MakeAlbums struct {
	Enabled bool
	Use     string
}
