# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this project adheres to [Semantic Versioning](https://semver.org/).

## 3.0.0
> This is a **major upgrade**, so it has some **non-backwards compatible changes**
### Added
- **Progress bar** when uploading files.
- Configuration, wo/ sensible data, is printed when debug is enabled. ([#270][i270]) 
- Configuration validation. The cli validates the configuration data at starting time.
- Information messages to bring more context at runtime. ([#260][i260]) 
### Changed
- `Jobs.MakeAlbums` configuration setting has changed to `Jobs.CreateAlbums`. Valid values are "Off", "folderName" and "folderPath".
- **Reduced the number of calls to the API when uploading files**. It's using less than 50% of calls than before.
- Move to `golang.org/x/term` from `golang.org/x/crypto/ssh/terminal`, due to deprecation.
- Some parts of the code has been refactored to make cleaner code and increase testability.
- `Jobs.Account` configuration setting has been changed to `Account`. Multiple Google Photos accounts are not supported. ([#231][i231]) 
- Bump `google-photos-api-client-go` from `v2.0.0` to `v2.1.3`. It improves performance. ([#259][i259])
- Bump `golangci-lint` from `1.30.0` to `1.34.1`.
### Deprecated
- `Jobs.MakeAlbums` configuration setting. Use `Jobs.CreateAlbums` instead.  See [configuration documentation][idocumentation].
- `Jobs.Account` configuration setting. Use `Account` instead. See [configuration documentation][idocumentation].
### Fixed
- '~' is not expanded when reading file. ([#268][i268])
### Removed
- Multiple Google Photos account support has been removed. You can use multiple configuration files in the same application folder. ([#231][i231]) 

[i270]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/270
[i268]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/268
[i260]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/260
[i259]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/259
[i231]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/231

## 2.0.1
### Changed
- Bump `google-photos-api-client-go` from `v2.0.0` to `v2.0.1`.
### Fixed
- Media item creation was failing when Google Photos was reporting errors on media creation. ([#262][i262])

[i262]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/262

## 2.0.0
> This is a **major upgrade**, and it has some **non-backwards compatible changes**:
> - `includePatterns` & `excludePatterns` configuration options has changed.
> - `includePatterns` has a new default (`_IMAGE_EXTENSIONS_`).
> - `uploadVideos` configuration option has been removed.
### Added
- Two new tagged patterns has been added: `_IMAGE_EXTENSIONS_`, matching [supported image file types](https://support.google.com/googleone/answer/6193313), and `_RAW_EXTENSIONS_`, matching [supported RAW file types](https://support.google.com/googleone/answer/6193313). ([#249][i249])
- Retries management. It's implementing exponential back-off with a maximum of 4 retries by default.  ([#253][i253]) 
### Changed
- `includePatterns` & `excludePatterns` configuration options has changed. It's using a new format, please review de [configuration documentation][idocumentation].
- By default, if `includePatterns` is empty, `_IMAGE_EXTENSIONS_` will be used. ([#249][i249])  
- Bump `google-photos-api-client-go` from `v2.0.0-beta-1` to `v2.0.0`.
### Fixed
- Symlinks are now supported when scanning a folder. ([#190][i190])
> **Note:** This application does not terminate if there are any non-terminating loops in the file structure.
- `includePatterns` works as expected, with a clearer (I hope so) format. ([#152][i152])  
### Removed
- Deprecated `uploadVideos` configuration option. It was deprecated in [v0.4.0](https://github.com/gphotosuploader/gphotos-uploader-cli/releases/tag/v0.4.0).

[i253]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/253                                          
[i249]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/249                                               
[i190]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/190
[i152]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/152

## 1.2.0
### Added
- `GPHOTOS_CLI_TOKENSTORE_KEY` env var could be used to read token store encryption key from. This allows you to run the CLI non-interactively. ([#224][i224])

[i224]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/224

## 1.1.2
### Fixed
- Fix homebrew tap creation. ([#233][i233])

[i233]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/233

## 1.1.1
### Changed
- Internal packages have been moved to `internal/` folder to discourage its usage.
- Bump `google-photos-api-client-go` from `v1.1.5` to `v2.0.0-beta-1`.
### Fixed
- Fix duplicated albums creation. ([#192][i192])

[i192]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/192

## 1.1.0
### Added
- Flag `--dry-run` to `push` command. It's useful to validate `includePatterns` and `excludePatterns` configuration. ([#216][i216])
### Changed
- The `init` command sets `deleteAfterUpload: false` as default value. ([#214][i214])
- CI has been moved to GitHub actions (previously was drone.io).
- Bump `gphotosuploader/googlemirror` from `v0.3.7` to `v0.4.0`.
- Bump `99designs/keyring` from `v1.1.2` to `v1.1.5`.

[i216]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/216
[i214]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/214

## 1.0.8
### Changed
- Update CI/CD minimum version to Go 1.13.
- Bump github.com/gphotosuploader/googlemirror to v0.3.7.

### Fixed
- File upload on album creation error. Thanks to [@albertvaka](https://github.com/albertvaka) ([#212][i212])

[i212]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/212

## 1.0.7
### Changed
- Bump github.com/int128/oauth2cli to v1.12.1 ([#206][i206])
- Bump golang.org/x/oauth2 to v0.0.0-20200107190931-bf48bf16ab8d ([#205][i205])

### Fixed
- Fix (temporary) OAuth broken process ([#181][i181])

[i181]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/181
[i206]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/206
[i205]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/205

## 1.0.6
### Changed
- Updated some dependencies

## 1.0.5
### Fixed
- Fix issue when installing CLI with `go get`. ([#183][i183])

[i183]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/182

## 1.0.4
### Fixed
- Fix `init` command error when it was used with root folders. ([#172][i172])

[i172]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/172

## 1.0.3
### Fixed
- Fix inconsistent use of `folderName` configuration option. ([#170][i170])

[i170]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/170

## 1.0.2
### Fixed
- Fix issue that hung the application when the `results` queue was full. This happened every time the number of files to upload was higher than 10x number of concurrent processes. ([#167][i167])

[i167]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/167

## 1.0.1
### Added
- New command `auth` to authenticate against Google Photos. It's useful to refresh authentication tokens. ([#125][i125])

[i125]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/125

## 1.0.0
> This is a **major upgrade**, and it has several **non-backwards compatible changes**. See more details below.
### Added
- New option for Album creation: `use: folderPath` will use the full folder path as Album name. See [config documentation][idocumentation]. ([#150][i150])
- New flags to control CLI verbosity: `--silent` suppress all logs except Fatal ones, `--debug` enable a lot of verbosity to logs.
- [CONTRIBUTING](CONTRIBUTING.md) guide line has been added.
- New Logger package to improve log readability.
### Changed
- **ATTENTION**: To upload items, you **must** run `gphotos-uploader-cli push`. The new `push` command substitutes `gphotos-uploader-cli`, that was working in previous versions.
- **ATTENTION**: New default config directory: `~/.gphotos-uploader-cli`. Copy your old configuration into the new folder or use `--config ~/.config/gphotos-uploader-cli` in every call.
- Default log verbosity is now `info` level, use `--debug` if you want more verbose output.
- [README](README.md) has been updated fixing some typos.
### Deprecated
- Once Go 1.13 has been published, previous Go 1.11 support is deprecated. This project will maintain compatibility with the last two major versions published.
- Configuration parameter `uploadVideos` has been deprecated in favor of `_ALL_VIDEO_FILES_` tagged pattern. See [configuration documentation][idocumentation] for details.
### Fixed
- Fix issue uploading photos without the correct file name. ([#158][i158])
- Fix issue uploading photos multiple times and ignoring others. ([#160][i160])

[i150]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/150
[i158]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/158
[i160]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/160

## 0.8.7
### Changed
- Update dependencies to newer versions: [gphotosuploader/google-photos-api-client-go](https://github.com/gphotosuploader/google-photos-api-client-go) v1.1.2 and [int128/oauth2cli](https://github.com/int128/oauth2cli) to v1.7.0.
- Update [golangci](https://github.com/golangci/golangci-lint) linter to version 1.20.0.
- `App` package to deal with specific application settings. 

## 0.8.6
### Changed
- Remove `build` from version. Now `version` has all the tag+build information.
### Fixed
- Fix duplicated album creation. ([#135][i135])

[i135]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/135

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
- `deleteAfterUpload` option has been reactivated, it was removed on v0.4.0. If you use this option in [config file][idocumentation] files will be deleted from local repository after being uploaded to Google Photos. (#25)
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
- Add two configuration options to include (`includePatterns`) and exclude (`excludePatterns`) files to be uploaded. See [configuration documentation][idocumentation] for details.

### Changed
- Reduce memory footprint simplifying objects overhead
- Configuration parameter `uploadVideos` is now using `includePatterns` and `excludePatterns` instead of detecting video format. **ATTENTION:** This option **will be deprecated** in the future in favor of `_ALL_VIDEO_FILES_` tagged pattern. See [configuration documentation][idocumentation] for details.

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
- Update [configuration documentation][idocumentation] to add `SecretsBackendType` (#83)
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

[idocumentation]: https://gphotosuploader.github.io/gphotos-uploader-cli/
