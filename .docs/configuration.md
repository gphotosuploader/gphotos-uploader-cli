# Configuration

Example configuration file (usually at `~/.config/gphotos-uploader-cli/config.hjson`:    

```hjson
{
  SecretsBackendType: auto,
  APIAppCredentials: {
    ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
    ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
  }
  jobs: [
    {
      account: youremail@gmail.com
      sourceFolder: ~/folder/to/upload
      makeAlbums: {
        enabled: true
        use: folderNames
      }
      deleteAfterUpload: false
      uploadVideos: true
    }
  ]
}
```
### SecretsBackendType
This option allows you to choose which backend will be used for secrets storage. You set `auto` to allow the application decide which one will be used given your environment.

Available options for secrets backends are:

```
"auto"              For auto backend selection
"secret-service"    For gnome-keyring support
"keychain"          For OS X keychain support
"kwallet"           For KDE Secrets Manager support
"wincred"           For Windows credentials support
"file"              For encrypted file support - needs interation to supply a symetric encryption key
"pass"              For Password Store support - needs user interation to supply a GPG pass key
```

Most of the times `auto` is the proper one. The application will try to use the existing backends in the order [defined by the library](https://github.com/99designs/keyring/blob/2c916c935b9f0286ed72c22a3ccddb491c01c620/keyring.go#L28):

```
// This order makes sure the OS-specific backends
// are picked over the more generic backends.
var backendOrder = []BackendType{
	// Windows
	WinCredBackend,
	// MacOS
	KeychainBackend,
	// Linux
	SecretServiceBackend,
	KWalletBackend,
	// General
	PassBackend,
	FileBackend,
}
```

## APIAppCredentials
The credentials that are provided are just example ones. 
Replace them with credentials you create at [Google Console](https://console.cloud.google.com/apis/api/photoslibrary.googleapis.com).

## jobs
List of folders to upload and upload options for each folder.

### `account`
Needs to be unique. It's the Google Account identity (e-mail address) where the files of this job are going to be uploaded.

### `sourceFolder`
The folder to upload from.
Must be an absolute path. Can expand the home folder tilde shorthand.

### `makeAlbums`
If makeAlbums.enabled set to true, use the last folder path component as album name.

### `deleteAfterUpload`
If set to true, media will be deleted from local disk after upload. 
To avoid data corruption, the uploader will double check that a the picture exists in your library and is visually similar to the one on the local disk before deleting any file.

### `uploadVideos`
If set to true, media items identified as a video will be uploaded. If false, they will be skipped. 