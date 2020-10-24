[![Go Report Card](https://goreportcard.com/badge/github.com/gphotosuploader/gphotos-uploader-cli)](https://goreportcard.com/report/github.com/gphotosuploader/gphotos-uploader-cli)
[![codebeat badge](https://codebeat.co/badges/9f3561ad-2838-456e-bc92-68988eeb376b)](https://codebeat.co/projects/github-com-gphotosuploader-gphotos-uploader-cli-master)
[![codecov](https://codecov.io/gh/gphotosuploader/gphotos-uploader-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/gphotosuploader/gphotos-uploader-cli)
[![GitHub release](https://img.shields.io/github/release/gphotosuploader/gphotos-uploader-cli.svg)](https://github.com/gphotosuploader/gphotos-uploader-cli/releases/latest)
[![GitHub](https://img.shields.io/github/license/gphotosuploader/gphotos-uploader-cli.svg)](LICENSE)
<!--- [![Snap Status](https://build.snapcraft.io/badge/gphotosuploader/gphotos-uploader-cli.svg)](https://build.snapcraft.io/user/gphotosuploader/gphotos-uploader-cli) --->

# Google Photos uploader CLI

Command line tool to mass upload media folders to your Google Photos account(s).    

While the official tool only supports Mac OS and Windows, this brings an uploader to Linux too. Lets you upload photos from, in theory, any OS for which you can compile a Go program.     

# Features

- **Customizable configuration**: via JSON-like config file.
- **Multiple Google accounts support**: upload your pictures to multiple accounts.
- **Filter files with patterns**: include/exclude files & folders using patterns (see [documentation](configuration.md)).
- **Resumable uploads**: Uploads can be resumed, saving time and bandwidth. 
- **File deletion after uploading**: Clean up local files after being uploaded.
- **Track already uploaded files**: uploads only new files to save bandwidth.
- **Secure**: logs you into Google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).

# Limitations
## Rate Limiting
Google Photos imposes a rate limit on all API clients. The quota limit for requests to the Library API is 10,000 requests per project per day. The quota limit for requests to access media bytes (by loading a photo or video from a base URL) is 75,000 requests per project per day.

## Photo storage and quality
All media items uploaded to Google Photos using the API [are stored in full resolution](https://support.google.com/photos/answer/6220791) at original quality. **They count toward the user’s storage**.

# Contributing
Help us make `gphotos-uploader-cli` the best tool for uploading your local pictures to Google Photos.

## Reporting Issues
If you find a bug while working with `gphotos-uploader-cli`, please [open an issue on GitHub](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?assignees=pacoorozco&labels=bug&template=bug_report.md) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feedback & Feature Requests
You are more than welcome to open issues in this project to:

- [give feedback](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?title=Feedback:)
- [suggest new features](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?labels=feature+request&template=feature_request.md)

## Contributing Code
This project is mainly written in Golang. If you want to contribute code, see [Contributing guide lines](CONTRIBUTING.md) for more information.

# License
 
 Use of this source code is governed by an MIT-style license that can be found in the LICENSE [MIT](LICENSE) file.
