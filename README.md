[![Build Status](https://travis-ci.org/nmrshll/gphotos-uploader-cli.svg?branch=master)](https://travis-ci.org/nmrshll/gphotos-uploader-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli)
[![GitHub release](https://img.shields.io/github/release/nmrshll/gphotos-uploader-cli.svg)](https://github.com/nmrshll/gphotos-uploader-cli/releases/latest)
[![GitHub](https://img.shields.io/github/license/nmrshll/gphotos-uploader-cli.svg)](LICENSE)
<!--- [![Snap Status](https://build.snapcraft.io/badge/nmrshll/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/nmrshll/gphotos-uploader-cli) --->

# Google Photos uploader CLI

Command line tool to mass upload media folders to your Google Photos account(s).    

While the official tool is only supports Mac OS and Windows, this brings an uploader to Linux too. Lets you upload photos from, in theory, any OS for which you can compile a Go program.     

# Features:

- specify folders to upload in config file
- upload to multiple google accounts
- include/exclude files & folders using patterns (see [documentation](.docs/configuration.md))
- ~~optionally delete objects after uploaḍ~~
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

Download the pre-compiled binaries from the [releases page](https://github.com/nmrshll/gphotos-uploader-cli/releases/latest) and copy to the desired location.

### Compiling from source

You can compile the source code in your system. **Go 1.11+** is required to compile this application:

```
$ git clone https://github.com/nmrshll/gphotos-uploader-cli
$ cd gphotos-uploader-cli
$ make build
```

Or you can use `go get` if you prefer it:

```
$ go get github.com/nmrshll/gphotos-uploader-cli
```

## Configure
First initialize the config file using this command:
```
$ gphotos-uploader-cli init
```

then modify it at `~/.config/gphotos-uploader-cli/config.hjson` to specify your configuration.

You can review the [documentation](.docs/configuration.md) to specify the folder to upload, add more Google Accounts and tune your configuration.

If you have problems, please open an [issue](https://github.com/nmrshll/gphotos-uploader-cli/issues). 

## Run

Once it's configured you can start uploading files in this way:
``` 
$ gphotos-uploader-cli
```    

# Contributing
Have improvement ideas or want to help ? Please start by opening an [issue](https://github.com/nmrshll/gphotos-uploader-cli/issues). 

# Related
- [google photos client library](https://github.com/nmrshll/google-photos-api-client-go)
- [oauth2-noserver](https://github.com/nmrshll/oauth2-noserver)

# License
 
 Use of this source code is governed by an MIT-style license that can be found in the LICENSE [MIT](LICENSE) file.
