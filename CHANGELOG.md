# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this project adheres to [Semantic Versioning](https://semver.org/).

### Unreleased
### Added
- Add two configuration options to include (`includePatterns`) and exclude (`excludePatterns`) files to be uploaded. See [configuration documentation](.docs/configuration.md) for details.

### Changed
- Reduce memory footprint simplifying objects overhead

### Fixed
- Fix folder typo on secrets backend storage

### Removed
- **ATTENTION:** `deleteAfterUpload` option has been temporarily removed. So no local file is removed by `gphotos-uplaoder-cli`.  See [issue #25](https://github.com/nmrshll/gphotos-uploader-cli/issues/25) for more details.

## 0.3.2
### Added
- Add `go get` installation method to [README](README.md)

### Changed
- Update `github.com/gphotosuploader/google-photos-api-client-go` to v1.0.1
- Update `github.com/gphotosuploader/googlemirror` to v0.3.2

### Fixed
- Update [configuration documentation](.docs/configuration.md) to add `SecretsBackendType` (#83)
- Typo on [README](README.md)

## 0.3.1
### Changed
- Move some dependencies to the new [gphotosuploader](https://github.com/gphotosuploader) organization
- `make test` is not as verbose as before. To make easier to see if there is an error
### Removed
- Removed some useless and local vendor files

## 0.3.0
### Added
- Support for [different secret backends](https://github.com/99designs/keyring). (#15, #41, #50, #51 and #52)
- Added test to completeuploads package
### Changed
- Document code in a more complete way
- Add `google.golang.org/api/photoslibrary` as vendor library, due to [Google's announcement](https://code-review.googlesource.com/c/google-api-go-client/+/39951) (#53)
- The `tokenstore` library has been modified to allow [new secrets backends](https://github.com/99designs/keyring)
### Fixed
- Fix installation instructions (#72)
### Removed
- `go get` installation method has been removed.

## 0.2.1 - 2019-06-18
### Fixed
- Fix [Go Report Card](https://goreportcard.com/report/github.com/nmrshll/gphotos-uploader-cli) issues

## 0.2.0 - 2019-06-18
### Added
- Support 5 concurrent uploads: **reduce API calls, speed things up** (#45)
- Added this [changelog](CHANGELOG.md) file

## 0.1.18 - 2019-06-16
### Changed
- Update github.com/h2non/filetype from v1.0.5 to v1.0.8 (#60)

### Fixed
- Fix mismatched type files (#38)

## 0.1.16 - 2019-06-16
### Fixed
- Fix goreleaser configuration (remove deprecated statement)
- Update [Getting started](README.md) documentation

### Removed
- Remove [snap](https://snapcraft.io/snaps) application publication (someone has stolen our app name)

## 0.1.11 - 2018-09-20
### Added
- [goreleaser](https://goreleaser.com/) will be on charge of publishing [binaries](https://github.com/nmrshll/gphotos-uploader-cli/releases) after new release is done

