[![Build Status](https://cloud.drone.io/api/badges/gphotosuploader/gphotos-uploader-cli/status.svg)](https://cloud.drone.io/gphotosuploader/gphotos-uploader-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/gphotosuploader/gphotos-uploader-cli)
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

by default configuration folder is `~/.config/gphotos-uploader-cli` but you can specify your own folder using `--config /my/config/dir`. Configuration is kept in the `config.hjson` file inside this folder.

You can review the [documentation](.docs/configuration.md) to specify the folder to upload, add more Google Accounts and tune your configuration.

If you have problems, please open an [issue](https://github.com/gphotosuploader/gphotos-uploader-cli/issues). 

## Run

Once it's configured you can start uploading files in this way:
``` 
$ gphotos-uploader-cli
```    

# Contributing
Have improvement ideas or want to help ? Please start by opening an [issue](https://github.com/gphotosuploader/gphotos-uploader-cli/issues). 


# License
 
 Use of this source code is governed by an MIT-style license that can be found in the LICENSE [MIT](LICENSE) file.
