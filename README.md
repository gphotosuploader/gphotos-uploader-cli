[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli)
<!--- [![Snap Status](https://build.snapcraft.io/badge/nmrshll/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/nmrshll/gphotos-uploader-cli) --->


# Google photos uploader CLI
Command line tool to mass upload media folders to your google photos account(s).    

While the official tool is only supports Mac OS and Windows, this brings an uploader to Linux too.

#### Features:
- specify folders to upload in config file
- optionally delete after upload
- upload to multiple google accounts
- security: logs you into google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).

# Quick start
##### Install using     
```
go get -u github.com/nmrshll/gphotos-uploader-cli
```    
##### Configure which folders to upload by modifying the file 
```
~/.config/gphotos-uploader-cli/config.hjson
```
([documentation](./docs/configuration.md))    
##### Run it with 
```
gphotos-uploader
```    

## Requirements
- Go 1.5+ for installation using `go install github.com/nmrshll/gphotos-uploader-cli`
- Mac OS or Linux

# Contributing
Have improvement ideas or want to help ? Please start by opening an [issue](https://github.com/nmrshll/gphotos-uploader-cli/issues)  

### Current plans
- [ ] add CI pipeline for testing / building / releasing deb/snap/homebrew/... packages (to drop the dependency on go for installing)
- [ ] add tests
- [ ] add CLI manual
- [ ] add electron app for front-end
- [ ] increase upload parallelism for speed

### Related
- [google photos client library](https://github.com/nmrshll/google-photos-api-client-go)
- [oauth2-noserver](https://github.com/nmrshll/oauth2-noserver)


#### License: [MIT](./.docs/LICENSE)
