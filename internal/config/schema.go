package config

// Config represents the content of configuration file.
// It defines the schema for Marshal and Unmarshal the data of the configuration file.
type Config struct {
	// APIAppCredentials represents Google Photos API credentials for OAuth.
	APIAppCredentials APIAppCredentials `json:"APIAppCredentials"`

	// Account is the Google Photos account to work with.
	Account string `json:"Account"`

	// SecretsBackendType is the type of backend to store secrets.
	SecretsBackendType string `json:"SecretsBackendType"`

	// Jobs are the source folders to work with.
	Jobs []FolderUploadJob `json:"Jobs"`
}

// APIAppCredentials represents Google Photos API credentials for OAuth.
type APIAppCredentials struct {
	// ClientID is the app identifier generated on the Google API console.
	ClientID string `json:"ClientID"`
	// ClientSecret is the secret key generated on the Google API console.
	ClientSecret string `json:"ClientSecret"`
}

// FolderUploadJob represents configuration for a folder to be uploaded
type FolderUploadJob struct {
	// SourceFolder is the folder containing the objects to be uploaded.
	SourceFolder string `json:"SourceFolder"`

	// CreateAlbums is the parameter to create albums on Google Photos.
	// Valid options are:
	// Off: Disable album creation (default).
	// folderPath: Creates album with the name based on full folder path.
	// folderName: Creates album with the name based on the folder name.
	CreateAlbums string `json:"CreateAlbums,omitempty"`

	// DeleteAfterUpload if it is true, the app will remove files after upload them.
	DeleteAfterUpload bool `json:"DeleteAfterUpload"`

	// IncludePatterns are the patterns to include files to work with.
	IncludePatterns []string `json:"IncludePatterns"`

	// ExcludePatterns are the patterns to exclude files.
	ExcludePatterns []string `json:"ExcludePatterns"`
}
