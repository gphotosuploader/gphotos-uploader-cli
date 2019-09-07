[![Build Status](https://cloud.drone.io/api/badges/gphotosuploader/gphotos-uploader-cli/status.svg)](https://cloud.drone.io/gphotosuploader/gphotos-uploader-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/gphotosuploader/gphotos-uploader-cli)
[![codebeat badge](https://codebeat.co/badges/9f3561ad-2838-456e-bc92-68988eeb376b)](https://codebeat.co/projects/github-com-gphotosuploader-gphotos-uploader-cli-master)
[![codecov](https://codecov.io/gh/gphotosuploader/gphotos-uploader-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/gphotosuploader/gphotos-uploader-cli)
[![GitHub release](https://img.shields.io/github/release/gphotosuploader/gphotos-uploader-cli.svg)](https://github.com/gphotosuploader/gphotos-uploader-cli/releases/latest)
[![GitHub](https://img.shields.io/github/license/gphotosuploader/gphotos-uploader-cli.svg)](LICENSE)
<!--- [![Snap Status](https://build.snapcraft.io/badge/gphotosuploader/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/gphotosuploader/gphotos-uploader-cli) --->

# Google Photos uploader CLI

Command line tool to mass upload media folders to your Google Photos account(s).    

While the official tool is only supports Mac OS and Windows, this brings an uploader to Linux too. Lets you upload photos from, in theory, any OS for which you can compile a Go program.     

# Features:

- specify folders to upload in config file
- upload to multiple google accounts
- include/exclude files & folders using patterns (see [documentation](.docs/configuration.md))
- resumable uploads
- optionally delete objects after uploaḍ
- security: logs you into google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).

# Getting started

## Install
You can install the pre-compiled binary (in several different ways) or compile from source.

Here are the steps for each of them:

### Install the pre-compiled binary

**homebrew tap** (only on macOS for now):
```
$ brew install gphotosuploader/tap/gphotos-uploader-cli
```

**manually**

Download the pre-compiled binaries from the [releases page](https://github.com/gphotosuploader/gphotos-uploader-cli/releases/latest) and copy to the desired location.

### Compiling from source

You can compile the source code in your system. **Go 1.11+** is required to compile this application:

```
$ git clone https://github.com/gphotosuploader/gphotos-uploader-cli
$ cd gphotos-uploader-cli
$ make build
```

Or you can use `go get` if you prefer it:

```
$ go get github.com/gphotosuploader/gphotos-uploader-cli
```

## Configure
First initialize the config file using this command:
```
$ gphotos-uploader-cli init
```

> Default configuration folder is `~/.config/gphotos-uploader-cli` but you can specify your own folder using `--config /my/config/dir`. Configuration is kept in the `config.hjson` file inside this folder.

You must review the [documentation](.docs/configuration.md) to specify your **Google Photos API credentials**, `APIAppCredentials`. You should tune your `jobs` configuration also.

## Run
Once it's configured you can start uploading files in this way:
``` 
$ gphotos-uploader-cli
```

### First time run
The first time you run `gphotos-uploader-cli`, after setting your configuration ([Google Photos API credentials](.docs/configuration.md#APIAppCredentials)), few manual steps are needed:

1. You should get an output like this one:

```
2019/09/07 17:47:11 Token has not been retrieved from token store: failed retrieving token from keyring
2019/09/07 17:47:11 Open http://localhost:40923
```

1. A browser will be opened at this point. Select the account where you wan to upload your files (the same you configured in the config file). You will see something like this:

![Google asking for Google Photos API credentials](.docs/images/ask_Google_Photos_API_credentials.png) 

1. After that, you should confirm that you trust on `gphotos-uploader-cli` to access to your Google Photos account, click on **Go to gphotos-uploader**:

![Google ask you to verify gphotos-upload-cli](.docs/images/ask_for_application_verification.png)

1. Finally Google will ask you to confirm permission Google Photos account:

![Google ask permission to your Google Photos account](.docs/images/ask_for_permission.png)

1. A green page is shown if all is fine

![Final confirmation, all was good](.docs/images/final_confirmation.png)

1. Go back to your terminal, you will see something like:

```
2019/09/07 17:47:23 Token expiration: 2019-09-07 18:47:23.588661948 +0200 CEST m=+3612.289671492
```

All auth configuration is in place.

# Contributing
Have improvement ideas or want to help ? Please start by opening an [issue](https://github.com/gphotosuploader/gphotos-uploader-cli/issues). 

# License
 
 Use of this source code is governed by an MIT-style license that can be found in the LICENSE [MIT](LICENSE) file.
