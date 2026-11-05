# Configuration

## Overview

The configuration file (`config.hjson`) controls how `gphotos-uploader-cli` works. By default, it is located in
`~/.gphotos-uploader-cli`, but you can specify a custom folder with `--config /my/config/dir`.

## Example Configuration

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
      Album: "template:%_directory%"
      DeleteAfterUpload: false
      IncludePatterns: [ "**/*.jpg", "**/*.png" ]
      ExcludePatterns: [ "**/ScreenShot*" ]
    }
  ]
}
```

## Configuration Options

### APIAppCredentials

OAuth 2.0 credentials for Google Photos API access.

**Setup Steps:**

1. Log in to your Google Account.
1. [Create a new project](https://console.cloud.google.com/projectcreate) in Google Cloud Platform and give it a name (
   example: _Google Photos Uploader_).
1. [Enable the Google Photos Library API](https://console.cloud.google.com/apis/library/photoslibrary.googleapis.com).
1. [Configure the OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent) by setting the
   application name (example: _gphotos-uploader-cli_) and then click the **Save** button on the bottom.
1. [Create OAuth credentials](https://console.cloud.google.com/apis/credentials) by clicking the **Create credentials →
   OAuth client ID** option, then choose **Desktop app** as the application type and give it a name (example:
   _gphotos-uploader-cli_).
1. Copy the **Client ID** and the **Client Secret** into your configuration file.

### Account

The Google Account (email) where files will be uploaded.

### SecretsBackendType

Choose where secrets (tokens) are stored. The `init` command sets it to `file`, which works everywhere but asks for a
passphrase (see [GPHOTOS_CLI_TOKENSTORE_KEY](#gphotos_cli_tokenstore_key)). On a desktop with a keyring available,
`auto` is usually more convenient.

Available options for secrets backend are:

| Option           | Description                                  |
|------------------|----------------------------------------------|
| `auto`           | Auto-select based on OS                      |
| `secret-service` | GNOME Keyring (Linux)                        |
| `keychain`       | macOS Keychain (macOS)                       |
| `kwallet`        | KDE Secrets Manager                          |
| `file`           | Encrypted file storage (requires passphrase) |

The application tries backends
in [this order](https://github.com/99designs/keyring/blob/2c916c935b9f0286ed72c22a3ccddb491c01c620/keyring.go#L28):
Keychain, SecretService, KWallet, File.

### Jobs

A list of upload jobs, each with its own options.

#### SourceFolder

Absolute path to the folder to upload. `~` is expanded to your home directory.

Symlinks are followed. Infinite loops are not detected.

#### Album

Controls album organization in Google Photos.

- **Omit**: Files are uploaded without an album. It's the default option.
- **Fixed Name**: `name:AlbumName` (uploads to the specific album _AlbumName_).
- **Template**: `template:...` (dynamic album names using placeholders, see below).

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

##### Fixed Album Name: `name:`

Specify `name:` followed by an album's name to upload objects to an album with the specified name. The album name in
Google Photos is not unique, so the first match will be used, or a new album will be created if none exists.

Example:

 ```hjson
  Album: name:fooBar
 ```

will create and upload objects to an album named `fooBar`:

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

##### Template-Based Album Names: `template:`

Use `template:` followed by a template string with placeholders to dynamically generate album names based on file
properties like date or folder names.

**Template Placeholders:**

| Placeholder           | Description                                                                         |
|-----------------------|-------------------------------------------------------------------------------------|
| `%_directory%`        | Name of the containing folder (same as the **deprecated** `auto:folderName` option) |
| `%_parent_directory%` | Name of the parent folder                                                           |
| `%_folderpath%`       | Path of the containing folder relative to `SourceFolder`, with `/` replaced by `_`. |
| `%_day%`              | Day of the month the file was last modified (in "DD" format).                       |
| `%_month%`            | Month the file was last modified (in "MM" format).                                  |
| `%_year%`             | Year the file was last modified (in "YYYY" format).                                 |
| `%_time%`             | Time the file was last modified (in "HH:MM:SS" 24-hour format).                     |
| `%_time_en%`          | Time the file was last modified (in "HH:MM:SS AM/PM" 12-hour format).               |

> Date and time placeholders use the local file's **modification time**, not its EXIF date. This only affects the
> album name: the date shown for each item in Google Photos is still taken from its EXIF data (or the upload date).

**Template Functions:**

| Function            | Description               |
|---------------------|---------------------------|
| $cutLeft(x,n)       | Remove first n characters |
| $cutRight(x,n)      | Remove last n characters  |
| $regex(x,expr,repl) | Regex replace             |
| $sentence(x)        | Sentence case             |
| $title(x)           | Title case                |
| $upper(x)           | Uppercase                 |
| $lower(x)           | Lowercase                 |

**Example:**

 ```hjson
  Album: template:%_directory% - %_month%.%_day%.$cutLeft(%_year%,2)
  ```

will calculate the album name based on the template for each file.

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

#### DeleteAfterUpload

If `true`, deletes local files after upload.

#### Including and Excluding files

Use the `IncludePatterns` and `ExcludePatterns` options to filter files. Each option is a list of one or more
patterns. Patterns are matched against paths relative to `SourceFolder`.

For example, to upload all _JPG and PNG files_ that are not named _*ScreenShots*_ you can configure it like this:

```hjson
IncludePatterns: [ "**/*.jpg", "**/*.png" ]
ExcludePatterns: [ "**/ScreenShot*" ]
```

Another example excluding a specific directory (and folders inside it):

```hjson
IncludePatterns: [ "_ALL_FILES_" ]
ExcludePatterns: [ "**/Temp/**" ]
```

> If `IncludePatterns` is empty, `_IMAGE_EXTENSIONS_` will be used.

**Special Patterns:**

| Pattern              | Description                                                                  |
|----------------------|------------------------------------------------------------------------------|
| `_ALL_FILES_`        | All files (**)                                                               |
| `_IMAGE_EXTENSIONS_` | [Supported image types](https://support.google.com/googleone/answer/6193313) |
| `_RAW_EXTENSIONS_`   | [Supported RAW types](https://support.google.com/googleone/answer/6193313)   |
| `_ALL_VIDEO_FILES_`  | [Supported video types](https://support.google.com/googleone/answer/6193313) |

**Pattern Syntax:**

- `*` matches any sequence of non-path-separators
- `**` matches any sequence, including path separators
- `?` matches a single non-path-separator
- `[class]` character classes
  - `[abc]` matches any single character within the set
  - `[a-z]` matches any single character in the range
  - `[^class]` matches any single character which does **not** match the class
- `{alt1,alt2}` alternatives
- Any character with a special meaning can be escaped with a backslash (`\`).

## Environment variables

### GPHOTOS_CLI_TOKENSTORE_KEY

Set this variable to provide the token store key when using `SecretsBackendType: file` (for headless runners).

```bash
GPHOTOS_CLI_TOKENSTORE_KEY=my-super-secret gphotos-uploader-cli push
```
