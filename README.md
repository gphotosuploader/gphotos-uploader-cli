[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli)
[![Snap Status](https://build.snapcraft.io/badge/nmrshll/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/nmrshll/gphotos-uploader-cli)

# Google photos uploader CLI
Command line tool to mass upload media folders to your google photos account(s).    
There is no official google photos desktop app for linux, this aims to fill this need.    
#### Features:
- specify folders to upload in config file
- optionally delete after upload
- uploading to multiple google accounts
- security: logs you into google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).

# Quick start
Install using     
`go get github.com/nmrshll/gphotos-uploader-cli`    
`go install github.com/nmrshll/gphotos-uploader-cli`    
Configure which folders to upload by modifying the file `~/.config/gphotos-uploader-cli/config.hjson` ([documentation](./docs/configuration.md))    
Run it with `gphotos-uploader`    

## Requirements
- Go 1.5+ for installation using `go install github.com/nmrshll/gphotos-uploader-cli` (running the)
- The linux keyring or macOS keychain (as of 2018-07)
- a unix-like filesystem (as of 2018-07)

# Contributing
Please submit an issue to discuss improvements before submitting a pull request.    

### Current issues
- [ ] add CI pipeline for testing / building / releasing deb/snap/homebrew/... packages (to drop the dependency on go for installing)
- [ ] add tests
- [ ] add CLI manual
- [ ] add electron app for front-end
- [ ] increase upload parallelism for speed

### Related
- [google photos client library](github.com/nmrshll/google-photos-api-client-go)
- [oauth2-noserver](github.com/nmrshll/oauth2ns)


#### License: [MIT](./.docs/LICENSE)
