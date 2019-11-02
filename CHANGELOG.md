# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this project adheres to [Semantic Versioning](https://semver.org/).

## Unreleased
### Added
- [CONTRIBUTING](CONTRIBUTING.md) guide line has been added.
- New Logger package to improve log readability.
### Changed
- [README](README.md) has been updated fixing some typos.
### Deprecated
- Once Go 1.13 has been published, previous Go 1.11 support is deprecated. This project will maintain compatibility with the last two major versions published.

## 0.8.7
### Changed
- Update dependencies to newer versions: [gphotosuploader/google-photos-api-client-go](https://github.com/gphotosuploader/google-photos-api-client-go) v1.1.2 and [int128/oauth2cli](https://github.com/int128/oauth2cli) to v1.7.0.
- Update [golangci](https://github.com/golangci/golangci-lint) linter to version 1.20.0.
- `App` package to deal with specific application settings. 

## 0.8.6
### Changed
- Remove `build` from version. Now `version` has all the tag+build information.
### Fixed
- Fix duplicated album creation. (#135)

## 0.8.5
### Fixed
- Fix exit code on clean executions. (#137)

## 0.8.4
### Fixed
- Fix panic when a unexpected error on media item creation was raised. (#110)
### Changed
- Update `gphotosuploader/google-photos-api-client-go` to v1.0.7.

## 0.8.3
### Fixed
- Fix never ending upload due to upstream bug [gphotosuploader/google-photos-api-client-go#19](https://github.com/gphotosuploader/google-photos-api-client-go/issues/19). (#129)
- Fix linter warning on `datastore/tokenstore/repository_keyring.go:100`.
- Fix [Codebeat](https://codebeat.co/projects/github-com-gphotosuploader-gphotos-uploader-cli-master) settings.
### Changed
- Updated `int128/oauth2cli` to version v1.5.0 and improve messages in logs. See [int128/oauth2cli#3](https://github.com/int128/oauth2cli/issues/3).

## 0.8.2
### Fixed
- Fix panic when a very big file was uploaded. It was solved in upstream [gphotosuploader/google-photos-api-client-go#17](https://github.com/gphotosuploader/google-photos-api-client-go/issues/17). (#127)

## 0.8.1
### Added
- Coverage reports in [codecov](https://codecov.io/gh/gphotosuploader/gphotos-uploader-cli/) service.

### Changed
- Updated `google-photos-api-client` to version v1.0.4 to help with broken album creation. (#19)

### Fixed
- Fix duplicated album creation due to concurrency problem. (#19)

## 0.8.0
### Added
- Uploads can be resumed. This will help uploading large files or when connection has fails. Thanks to @pdecat.

## 0.7.2
### Fixed
- Fix token storing when expired token has been refreshed. See comments on #107.
- Fix typos in CHANGELOG.

## 0.7.1 
### Added
- Google Auth expired token refresh. Once token is expired, `gphotos-uploader-cli` will try to refresh the token without user intervention. **NOTE**: First time you use this version, you should re-authenticate in order to get the token that allows token refresh. (#103)
- Add `--config` flag to specify the folder where configuration is kept. (#104)
### Changed
- Moved CI/CD platform from Travis to [Drone.io](https://cloud.drone.io/gphotosuploader/gphotos-uploader-cli). It has reduce the time to CI by a half.

## 0.6.0
### Added
- `deleteAfterUpload` option has been reactivated, it was removed on v0.4.0. If you use this option in [config file](.docs/configuration.md) files will be deleted from local repository after being uploaded to Google Photos. (#25)
### Changed
- This repository has transferred to [GPhotos Uploaders organization](https://github.com/gphotosuploader), so all imports has been updated to the new organization's URL.
### Removed
- Removed some useless log lines. There are still too much.

## 0.5.0
### Changed
- Fix issue #97 "New gnome keyring store created on each launch". To solve this issue a new `serviceName` has been changed. **NOTE**: Once you use this version, a new Gnome keyring will be created, so credentials should be supplied again. (#97) 

## 0.4.2
### Fixed
- Fix CI release pipeline to fix application version (#94). Last version was still broken on CI.

## 0.4.1
### Added
- Add Homebrew tap to allow users to install `gphotos-uploader-cli` using Homebrew. See [install](README.md) section.
 
### Fixed
- Fix CI release pipeline to fix application version (#94)

## 0.4.0
### Added
- Add two configuration options to include (`includePatterns`) and exclude (`excludePatterns`) files to be uploaded. See [configuration documentation](.docs/configuration.md) for details.

### Changed
- Reduce memory footprint simplifying objects overhead
- Configuration parameter `uploadVideos` is now using `includePatterns` and `excludePatterns` instead of detecting video format. **ATTENTION:** This option **will be deprecated** in the future in favor of `_ALL_VIDEO_FILES_` tagged pattern. See [configuration documentation](.docs/configuration.md) for details.

### Fixed
- Fix folder path typo on secrets backend storage

### Removed
- **ATTENTION:** `deleteAfterUpload` option has been temporarily removed. So no local file is removed by `gphotos-uplaoder-cli`.  See [issue #25](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/25) for more details.

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
- Fix [Go Report Card](https://goreportcard.com/report/github.com/gphotosuploader/gphotos-uploader-cli) issues

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
- [goreleaser](https://goreleaser.com/) will be on charge of publishing [binaries](https://github.com/gphotosuploader/gphotos-uploader-cli/releases) after new release is done

