package config

// Config represents the content of the configuration file.
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

	// Album is the album where objects will be uploaded.
	// If the Album option is not set, the objects will not be associated with an album in Google Photos.
	//
	// These are the valid values: "name:", "auto:", "template".
	//   "name:" : Followed by the album name in Google Photos (album names are not unique, so the first to match
	//             will be selected)
	//   "template": Followed by a template string that can contain the following predefine tokens and functions:
	//              Tokens:
	//                 %_folderpath% - full path of the folder containing the file.
	//                 %_directory% - name of the folder containing the file.
	//                 %_parent_directory% - Replaced with the name of the parent folder of the file.
	//                 %_day% - day of the month the file was created (in "DD" format).
	//                 %_month% -month the file was created (in "MM" format).
	//                 %_year% - year the file was created (in "YYYY" format).
	//                 %_time% - time the file was created (in "HH:MM:SS" 24-hour format).
	//                 %_time_en% - time the file was created (in "HH:MM:SS AM/PM" 12-hour format).
	//              Functions:
	//                 $lower(string) - converts the string to lowercase.
	//                 $upper(string) - converts the string to uppercase.
	//                 $sentence(string) - converts the string to sentence case.
	//                 $title(string) - converts the string to title case.
	//                 $regex(string, regex, replacement) - replaces the string with the regex replacement.
	//                 $cutLeft(string, length) - cuts the string from the left.
	//                 $cutRight(string, length) - cuts the string from the right.
	//
	//              Example: "template:%_directory% - %_month%.%_day%.$cutLeft(%_year%,2)"

	Album string `json:"Album,omitempty"`

	// CreateAlbums exists to notice users about its deprecation. It should not be used in favor of the Album option.
	CreateAlbums string `json:"CreateAlbums,omitempty"`

	// DeleteAfterUpload if it is true, the app will remove files after upload them.
	DeleteAfterUpload bool `json:"DeleteAfterUpload"`

	// IncludePatterns are the patterns to include files to work with.
	IncludePatterns []string `json:"IncludePatterns"`

	// ExcludePatterns are the patterns to exclude files.
	ExcludePatterns []string `json:"ExcludePatterns"`
}
