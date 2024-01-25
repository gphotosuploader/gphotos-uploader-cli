# Configuration

## Configuration options

> The configuration is kept in the file `config.hjson` inside the configuration folder. You can specify your own folder using `--config /my/config/dir` otherwise default configuration folder is `~/.gphotos-uploader-cli`.

Example configuration file:    

```hjson
{
  APIAppCredentials:
  {
    ClientID: YOUR_APP_CLIENT_ID
    ClientSecret: YOUR_APP_CLIENT_SECRET
  }
  Account: YOUR_GOOGLE_PHOTOS_ACCOUNT
  SecretsBackendType: file
  Jobs:
  [
    {
      SourceFolder: YOUR_FOLDER_PATH
      Album: "auto:folderName"
      DeleteAfterUpload: false
      IncludePatterns: [ "**/*.jpg", "**/*.png" ]
      ExcludePatterns: [ "**/ScreenShot*" ]
    }
  ]
}
```

## APIAppCredentials <!-- {docsify-ignore} -->

Given that `gphotos-uploader-cli` uses OAuth 2 to access Google APIs, authentication is a bit tricky and involves a few manual steps. Please follow the guide below carefully, to give `gphotos-uploader-cli` the required access to your Google Photos account.

Before you can use `gphotos-uploader-cli`, you must enable the Photos Library API and request an OAuth 2.0 Client ID.

1. Make sure you're logged in into the Google Account where your photos should be uploaded to.
1. Start by [creating a new project](https://console.cloud.google.com/projectcreate) in Google Cloud Platform and give it a name (example: _Google Photos Uploader_).
1. Enable the [Google Photos Library API](https://console.cloud.google.com/apis/library/photoslibrary.googleapis.com) by clicking the <kbd>ENABLE</kbd> button.
1. Configure the [OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent) by setting the application name (example: _gphotos-uploader-cli_) and then click the <kbd>Save</kbd> button on the bottom.
1. Create [credentials](https://console.cloud.google.com/apis/credentials) by clicking the **Create credentials → OAuth client ID** option, then pick **Desktop app** as the application type and give it a name (example: _gphotos-uploader-cli_).
1. Copy the **Client ID** and the **Client Secret** and keep them ready to use in the next step.
1. Open the *config file* and set both the `ClientID` and `ClientSecret` options to the ones generated on the previous step.

## Account
It's the Google Account identity (e-mail address) where the files are going to be uploaded.

### SecretsBackendType <!-- {docsify-ignore} -->
This option allows you to choose which backend will be used for secret storage. You set `auto` to allow the application to decide which one will be used given your environment.

Available options for secrets backend are:

```
"auto"              For auto backend selection
"secret-service"    For gnome-keyring support
"keychain"          For OS X keychain support
"kwallet"           For KDE Secrets Manager support
"file"              For encrypted file support - needs interaction to supply a symmetric encryption key
```

Most of the time `auto` is the proper one. The application will try to use the existing backends in the order [defined by the library](https://github.com/99designs/keyring/blob/2c916c935b9f0286ed72c22a3ccddb491c01c620/keyring.go#L28):

```
// This order makes sure the OS-specific backends
// are picked over the more generic backends.
var backendOrder = []BackendType{
	// MacOS
	KeychainBackend,
	// Linux
	SecretServiceBackend,
	KWalletBackend,
	// General
	FileBackend,
}
```

## Jobs <!-- {docsify-ignore} -->
List of folders to upload and upload options for each folder.

### SourceFolder
The folder to upload from. Must be an absolute path. Can expand the home folder tilde shorthand `~`.
> The application will follow any symlink it finds, it does not terminate if there are any non-terminating loops in the file structure.

### Album
It controls how uploaded files will be organized into albums in Google Photos.

Given the local tree of folders and files:

```shell
/home/my-user/pictures
└── upload
    ├── album-1
    │   ├── image-album1-01.jpg
    │   └── image-album1-02.jpeg
    ├── album-2
    │   ├── image-album2-01.jpg
    │   └── image-album2-02.jpg
    └── album-3
        ├── image-album3-01.jpg
        ├── image-album3-02.jpg
        └── image-album3-03.jpg
```

These are several options: `name:`, `auto:`, `template:`:

#### Fixed name: `name:`

The `name:` option followed by an album's name, will upload objects to an album with the specified name. 

The album name in Google Photos is not unique, so the first to match to the name will be selected.

Setting `Album: name:fooBar` will create and upload objects to an album named `fooBar`:

```shell
Google Photos
└── fooBar
    ├── image-album1-01.jpg
    ├── image-album1-02.jpeg
    ├── image-album2-01.jpg
    ├── image-album2-02.jpg
    ├── image-album3-01.jpg
    ├── image-album3-02.jpg
    └── image-album3-03.jpg
```

#### Calculated name from a file path: `auto:`

##### From parent folder: `auto:folderName`

Setting `auto:folderName` and `SourceFolder: /home/my-user/pictures` will use the name of the folder (within `SourceFolder`), where the item is uploaded from, to set the album name.

```shell
Google Photos
├── album-1
│   ├── image-album1-01.jpg
│   ├── image-album1-02.jpeg
├── album-2
│   ├── image-album2-01.jpg
│   └── image-album2-02.jpg
└── album-3
    ├── image-album3-01.jpg
    ├── image-album3-02.jpg
    └── image-album3-03.jpg
```
##### From full path: `auto:folderPath` 

Setting `auto:folderPath` and `SourceFolder: /home/my-user/pictures` will use the full path of the folder (relative to `SourceFolder`), where the item is uploaded from, to set the album name.

```shell
Google Photos
├── upload_album-1
│   ├── image-album1-01.jpg
│   ├── image-album1-02.jpeg
├── upload_album-2
│   ├── image-album2-01.jpg
│   └── image-album2-02.jpg
└── upload_album-3
    ├── image-album3-01.jpg
    ├── image-album3-02.jpg
    └── image-album3-03.jpg
```

#### Customized template: `template:`

Using `template:` followed by a template string that can contain the following predefined tokens and functions:

##### Tokens

| Placeholder         | Description                                                              |
|---------------------|--------------------------------------------------------------------------|
| %_directory%        | Name of the folder containing the file (same as `auto:folderName`).      |
| %_parent_directory% | Name of the grandparent folder of the file.                              |
| %_folderpath%       | Full path of the folder containing the file (same as `auto:folderPath`). |
| %_day%              | Day of the month the file was created (in "DD" format).                  |
| %_month%            | Month the file was created (in "MM" format).                             |
| %_year%             | Year the file was created (in "YYYY" format).                            |
| %_time%             | Time the file was created (in "HH:MM:SS" 24-hour format).                |
| %_time_en%          | Time the file was created (in "HH:MM:SS AM/PM" 12-hour format).          |

##### Functions

| Placeholder         | Description                                                                                                                                                                                                                                                           |
|---------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| $cutLeft(x,n)       | Removes the first n characters of string x and returns the result.                                                                                                                                                                                                    |
| $cutRight(x,n)      | Removes the last n characters of string x and returns the result.                                                                                                                                                                                                     |
| $regex(x,expr,repl) | Replaces the pattern specified by the regular expression expr in the string x by repl. The fourth optional parameter enables ignore case (1) or disables the ignore case setting (0). Please note that you have to escape comma and other special characters in expr. |
| $sentence(x)        | Converts the given string to sentence case.                                                                                                                                                                                                                           |
| $title(x)           | Converts the given string to title case.                                                                                                                                                                                                                              |
| $upper(x)           | Converts the given string to upper case.                                                                                                                                                                                                                              |
| $lower(x)           | Converts the given string to lower case.                                                                                                                                                                                                                              |

##### Examples

Setting `template:%_directory% - %_month%.%_day%.$cutLeft(%_year%,2)` will calculate the album name based on the template for each file.

```shell
Google Photos
├── album-1 - 11.21.23
│   ├── image-album1-01.jpg
│   └── image-album1-02.jpeg
├── album-2 - 10.20.23
│   └── image-album2-01.jpg
├── album-2 - 10.22.23
│   └── image-album2-02.jpg
└── album-3 - 09.12.23
    ├── image-album3-01.jpg
    ├── image-album3-02.jpg
    └── image-album3-03.jpg
```

### DeleteAfterUpload
If set to true, media will be deleted from the local disk after completing the upload. 

## Including and Excluding files
You can include and exclude files by specifying the `includePatterns` and `excludePatterns` options. You can add one or more patterns separated by commas `,`. These patterns are always applied to `sourceFolder`.

For example, to upload all _JPG and PNG files_ that are not named _*ScreenShots*_ you can configure it like this:
```
includePatterns: [ "**/*.jpg", "**/*.png" ]
excludePatterns: [ "**/ScreenShot*" ]
```

Another example excluding a specific directory (and folders inside it):
```
includePatterns: [ "_ALL_FILES_" ]
excludePatterns: [ "**/Temp/**" ]
```

> If `includePatterns` is empty, `_IMAGE_EXTENSIONS_` will be used.

### Patterns
Supports the following special terms in the patterns:

Special Terms | Meaning
------------- | -------
`*`           | matches any sequence of non-path-separators
`**`          | matches any sequence of characters, including path separators
`?`           | matches any single non-path-separator character
`[class]`     | matches any single non-path-separator character against a class of characters ([see below](#character-classes))
`{alt1,...}`  | matches a sequence of characters if one of the comma-separated alternatives matches

Any character with a special meaning can be escaped with a backslash (`\`).

#### Character Classes

Character classes support the following:

Class      | Meaning
---------- | -------
`[abc]`    | matches any single character within the set
`[a-z]`    | matches any single character in the range
`[^class]` | matches any single character which does *not* match the class

#### Tagged patterns
There are some common patterns that have been tagged, you can use them to simplify your configuration.

> Tagged patterns matches file extensions case insensitively.

* `_ALL_FILES_`: Matches all files, is the same as using `**`. 
* `_IMAGE_EXTENSIONS_`: Matches [Google Photos supported image file types](https://support.google.com/googleone/answer/6193313) and it includes: `jpg, jpeg, png, webp, gif` file extensions case in-sensitively.
* `_RAW_EXTENSIONS_`: Matches [Google Photos supported RAW file types](https://support.google.com/googleone/answer/6193313) and it includes `arw, srf, sr2, crw, cr2, cr3, dng, nef, nrw, orf, raf, raw, rw2` file extensions case in-sensitively.
* `_ALL_VIDEO_FILES_`: Matches [Google Photos supported video file types](https://support.google.com/googleone/answer/6193313) and it includes `mpg, mod, mmv, tod, wmv, asf, avi, divx, mov, m4v, 3gp, 3g2, mp4, m2t, m2ts, mts, mkv` file extensions case in-sensitively.

## Environment variables

### GPHOTOS_CLI_TOKENSTORE_KEY

This variable is used to read the token store key for opening the secrets storage. It works when `SecretsBackendType: file` and it is intended to be used by headless runners.

```bash
GPHOTOS_CLI_TOKENSTORE_KEY=my-super-secret gphotos-uploader-cli push
```
