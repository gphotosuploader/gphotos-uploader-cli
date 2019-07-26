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
      includePatterns: [ "*.jpg", "*.png" ]
      excludePatterns: [ "*ScreenShot*" ]
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
> **ATTENTION:** For v0.4.0+ this option has been disabled due to several problems detecting similarities between uploaded and local files. **So no file will be removed from local storage by this CLI** See [issue #45](https://github.com/nmrshll/gphotos-uploader-cli/issues/25) for more details.
  
~~If set to true, media will be deleted from local disk after upload. 
To avoid data corruption, the uploader will double check that a the picture exists in your library and is visually similar to the one on the local disk before deleting any file.~~

### `uploadVideos`
> **DEPRECATION NOTICE:** This option will be deprecated in the future in favor of `_ALL_VIDEO_FILES_` tagged pattern used in `includePatterns` or `excludePatterns`.

If set to `true`, video media items will be included (uploaded). If `false`, video files will be excluded.

If you want to upload video files (`uploadVideos: true`), you can use:
```
includePatterns: [ "_ALL_VIDEO_FILES_" ]
excludePatterns: []
```

If you don't want to upload video files (`uploadVideos: false`), you can use:
```
includePatterns: [ "_ALL_FILES_" ]
excludePatterns: [ "_ALL_VIDEO_FILES_" ]
```

**NOTE:** It means that as far as `uploadVideos` options is present, video files are always included or excluded,

## Including and Excluding files
You can include and exclude files by specifying the `includePatterns` and `excludePatterns` options. You can add one or more patterns separated by commas `,`. These patterns are always applied to `sourceFolder`.

For example, to upload all _JPG and PNG files_ that are not named _*ScreenShots*_ you can configure it like this:
```
includePatterns: [ "*.jpg", "*.png" ]
excludePatterns: [ "*ScreenShot*" ]
```

### Patterns
Patterns use [filepath.Glob](https://golang.org/pkg/path/filepath/#Glob) internally, see [filepath.Match](https://golang.org/pkg/path/filepath/#Match) for syntax. 

Regular wildcards cannot be used to match over the directory separator `/`. For example: `b*ash` matches `/dir/bash` but does not match `/dirb/ash`.

For this, the special wildcard `**` can be used to match arbitrary sub-directories: The pattern `foo/**/bar` matches:

* `/dir1/foo/dir2/bar/file`
* `/foo/bar/file`
* `/tmp/foo/bar`

#### Tagged patterns
There are some common patterns that has been tagged, you can use them to simplify your configuration.

* `_ALL_FILES_`: Matches all files, is the same as using `*`. 
* `_ALL_VIDEO_FILES_`: Matches all video file extensions supported by Google Photos.
> Supported video extensions are sourced by [Google Photos support](https://support.google.com/googleone/answer/6193313) and it includes:
> .mpg, .mod, .mmv, .tod, .wmv, .asf, .avi, .divx, .mov, .m4v, .3gp, .3g2, .mp4, .m2t, .m2ts, .mts, and .mkv files.
