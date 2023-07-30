package configuration

// Config represents the content of the configuration file.
// It defines the schema for Marshal and Unmarshal the data of the configuration file.
type Config struct {
	// Auth represents the auth configuration.
	Auth AuthConfiguration `mapstructure:"auth"`

	// Folders are the source folders to work with.
	Folders []FolderUpload `mapstructure:"folders"`
}

// AuthConfiguration represents Google Photos API credentials for OAuth.
type AuthConfiguration struct {
	// ClientID is the app identifier generated on the Google API console.
	ClientID string `mapstructure:"client_id"`
	// ClientSecret is the secret key generated on the Google API console.
	ClientSecret string `mapstructure:"client_secret"`

	// Account is the Google Photos account to work with.
	Account string `mapstructure:"account"`

	// SecretsBackendType is the type of backend to store secrets.
	SecretsBackendType string `mapstructure:"secrets_type"`
}

// FolderUpload represents configuration for a folder to be uploaded
type FolderUpload struct {
	// Path is the folder containing the objects to be uploaded.
	Path string `mapstructure:"path"`

	// CreateAlbums is the parameter to create albums on Google Photos.
	// Valid options are:
	// Off: Disable album creation (default).
	// folderPath: Creates album with the name based on full folder path.
	// folderName: Creates album with the name based on the folder name.
	CreateAlbums string `mapstructure:"create_albums"`

	// DeleteAfterUpload if it is true, the app will remove files after upload them.
	DeleteAfterUpload bool `mapstructure:"delete_after_upload"`

	// Include are the patterns to include files to work with.
	Include []string `mapstructure:"include"`

	// Exclude are the patterns to exclude files.
	Exclude []string `mapstructure:"exclude"`
}
