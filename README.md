[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli)
[![GitHub release](https://img.shields.io/github/release/nmrshll/gphotos-uploader-cli.svg)](https://github.com/nmrshll/gphotos-uploader-cli/releases/latest)
[![GitHub](https://img.shields.io/github/license/nmrshll/gphotos-uploader-cli.svg)](LICENSE)
<!--- [![Snap Status](https://build.snapcraft.io/badge/nmrshll/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/nmrshll/gphotos-uploader-cli) --->

# Google Photos uploader CLI

Command line tool to mass upload media folders to your Google Photos account(s).    

While the official tool is only supports Mac OS and Windows, this brings an uploader to Linux too. Lets you upload photos from, in theory, any OS for which you can compile a Go program.     

# Features:

- specify folders to upload in config file
- optionally delete obejcts after upload
- upload to multiple google accounts
- security: logs you into google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).

# Quick start
There are two ways to install this command line:

### Downloading a binary (recommended)

You can download [latest released package](https://github.com/nmrshll/gphotos-uploader-cli/releases/latest) for your preferred OS.

[![GitHub release](https://img.shields.io/github/release/nmrshll/gphotos-uploader-cli.svg?style=for-the-badge)](https://github.com/nmrshll/gphotos-uploader-cli/releases/latest)

### Compiling in your system

You can compile the source code in your system. **Go 1.11+** is required to compile this application:

Clone this repository and use `make build` to compile it

```
$ git clone https://github.com/nmrshll/gphotos-uploader-cli
$ cd gphotos-uploader-cli
$ make build
```

**NOTE**: It's not possible to install this application using `go get` due to the [use of a third party library as vendor](https://github.com/golang/go/issues/30354).

## Configuring this application
First initialize the config file using this command:
```
$ gphotos-uploader-cli init
```

then modify it at `~/.config/gphotos-uploader-cli/config.hjson` to specify your configuration.

You can review the [documentation](.docs/configuration.md) to specify the folder to upload, add more Google Accounts and tune your configuration.

If you have problems, please take a look to [troubleshooting](.docs/installation-troubleshooting.md) guide.

## Run it with 

Once it's configured you can call the command in this way:
``` 
$ gphotos-uploader-cli
```    

# Contributing
Have improvement ideas or want to help ? Please start by opening an [issue](https://github.com/nmrshll/gphotos-uploader-cli/issues). 

## Current plans
- [ ] add CI pipeline for testing / building / releasing deb/snap/homebrew/... packages (to drop the dependency on go for installing)
- [ ] add tests
- [ ] add CLI manual
- [ ] add electron app for front-end
- [x] increase upload parallelism for speed

# Related
- [google photos client library](https://github.com/nmrshll/google-photos-api-client-go)
- [oauth2-noserver](https://github.com/nmrshll/oauth2-noserver)

# License
 
 Use of this source code is governed by an MIT-style license that can be found in the LICENSE [MIT](LICENSE) file.
