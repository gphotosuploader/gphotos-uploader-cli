[![Go Report Card](https://goreportcard.com/badge/github.com/nmrshll/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli)

# Google photos uploader CLI
Command line tool to upload large amounts of media to your google photos account(s)

## Quick start
Install using `go install github.com/nmrshll/gphotos-uploader-cli`
Configure by modifying the file `~/.config/gphotos-uploader-cli/config.hjson`
Run it with `gphotos-uploader`

# Why this tool ?
Google released apps for google photos for PC, Mac, Android and iOS, but Linux was ignored.
This tool intends to be a cross platform google photos app for uploading media.
Right now it's a Linux and Mac uploader CLI, I hope to add a graphical interface to it in the future.

This tool relies on the [google photos client library](github.com/nmrshll/google-photos-api-client-go)


## Features
- upload large amounts of files
- the program is safe to interrupt and restart later
- delete local files after upload is successful
- supports uploading to several google accounts
- Security
    - Logging in is done using OAuth (using lib [oauth2-noserver](github.com/nmrshll/oauth2ns)), which means you never have to trust my code with your google credentials
    - your OAuth token is stored inside the keyring/keychain for more security

## Requirements
- Go 1.5+ for installation using `go install github.com/nmrshll/gphotos-uploader-cli` (running the)
- The linux keyring or macOS keychain (as of 2018-07)
- a unix-like filesystem (as of 2018-07)

## Contributing
Please submit an issue to discuss improvements before submitting a pull request.

### Todo
- add CI pipeline for testing / building / releasing deb/snap/homebrew/... packages (to drop the dependency on go for installing)
- add tests
- add CLI manual
- add electron app for front-end
- increase upload parallelism for speed


## License
MIT License. Use for any purpose, including commercial.
