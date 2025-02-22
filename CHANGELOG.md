# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this project adheres to [Semantic Versioning](https://semver.org/).

## 5.0.1
### Added
- Support for go v1.24

### Changed
- Bump several dependencies

### Removed
- Support for go v1.22

## 5.0.0
### Added
- Google Photos scopes have changed, some of them have been removed. The CLI will use the new ones. ([#474][i474])

### Changed
- Bump `golang.org/x/text` to 0.19.0 ([#479][i479])
- Bump `golang.org/x/term` to 0.25.0 ([#478][i478])
- Bump `github.com/schollz/progressbar/v3` to 3.16.1 ([#477][i477])

### Removed
- The deprecated `Album: auto:folderName` and `Album: auto:folderPath` options have been removed. Use the `Album: template:%_directory%` and `Album: template:%_folderpath%` options instead.
- The deprecated `Jobs: CreateAlbums` option has been removed. Use the `Jobs: Album` option instead.

[i474]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/474
[i479]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/479
[i478]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/478
[i477]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/477

## 4.6.0
### Added
- Support for the latest published Go version (1.23). This project will maintain compatibility with the latest **two major versions** published.

### Changed
- Bump `gphotosuploader/google-photos-api-client-go/v3` to version 3.0.6
- Bump `schollz/progressbar/v3` to version 3.15.0
- Bump `spf13/cobra` to version 1.8.1


## 4.5.0
### Added
- Support for the latest published Go version (1.22). This project will maintain compatibility with the latest **two major versions** published.

### Changed
- Bump `golang.org/x/oauth2` to version 0.17.0
- Bump `gphotosuploader/google-photos-api-client-go/v3` to version 3.0.5

### Deprecated
- The `auto:folderName` and `auto:folderPath` options are deprecated in favor of the `template:%_directory%` and `template:%_folderpath%` options. See [documentation](https://gphotosuploader.github.io/gphotos-uploader-cli/#/configuration?id=customized-template-template).


## 4.4.0
### Added
- Option to customize Album names by introducing `template`. Thanks to [@WACKYprog](https://github.com/WACKYprog) ([#431][i431])

[i431]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/431

## 4.3.0
### Added
- Option to bind the HTTP server to address other than local ([#426][i426])

### Changed
- Bump `golang.org/x/term` to version 0.16.0
- Bump `golang.org/x/sync` to version 0.6.0
- Bump `github.com/dvsekhvalnov/jose2go` to version 1.6.0
- Bump `github.com/hjson/hjson-go/v4` to version 4.4.0
- Bump `golang.org/x/oauth2` to version 0.16.0
- Bump `github.com/spf13/afero` to version 1.11.0
- Bump `github.com/gphotosuploader/google-photos-api-client-go/v3` to version 3.0.4

[i426]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/426

## 4.2.1
### Fixed
- Fix versioning on releases since the publication of 4.x version. ([#413][i413])
- Small typos in messages. Thanks, [@tbm](https://github.com/tbm) ([#414][i414])

### Changed
- Bump `github.com/schollz/progressbar/v3` from 3.13.1 to 3.14.1 ([#411][i411])
- Bump `golang.org/x/oauth2` from 0.13.0 to 0.14.0 ([#409][i409])
- Bump `golang.org/x/sync` from 0.4.0 to 0.5.0 ([#408][i408])
- Bump `github.com/spf13/cobra` from 1.7.0 to 1.8.0 ([#407][i407])

[i413]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/413
[i414]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/414
[i411]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/411
[i409]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/409
[i408]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/408
[i407]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/407

## 4.2.0
### Added
- New parameter `--redirect-url-hostname` to the `auth` command in order to set the URL to use after the Google Photos authentication ([#402][i402])

### Changed
- Bump `github.com/hjson/hjson-go/v4` from 4.3.0 to 4.3.1 ([#400][i400])


[i400]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/400
[i402]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/402

## 4.1.1
### Fixed
- Uploads are showing the data of upload instead of the filename in the Google Photos UI.  Thanks [@mikebilly](https://github.com/mikebilly) ([#398][i398])

[i398]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/398

## 4.1.0
### Added
- Flag `--port` to configure the port where the authentication server will listen to when using the `auth` command ([#370][i370])
- **New command to reset the already uploaded file tracker** (`reset file-tracker`), which removes the internal database ([#182][i182])
- New **`Album` option in Job's configuration** which allows to set a fixed album's name to upload objects to. ([#393][i393])

### Deprecated
- The `CreateAlbums` option in Job's configuration is deprecated in favor of a new `Album` option.

[i370]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/370
[i182]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/182
[i393]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/393

## 4.0.0
### Added
- Support for the latest published Go version (1.21). This project will maintain compatibility with the latest **two major versions** published.
- Implement a cache to reduce the number of requests to Google Photos API and reduce the risk of being quota limited.
- Implement **a new command to list albums** (`list albums`) created by this CLI.
- Implement **a new command to list media items** (`list media-items`) uploaded by this CLI. It offers the possibility of filtering by album.
- Progress bars to provide feedback to users on very long transactions. 

### Changed
- Bump `github.com/sirupsen/logrus` from 1.9.0 to 1.9.3 ([#378][i378])
- Bump `github.com/spf13/afero` from 1.9.5 to 1.10.0 ([#379][i379])
- Bump `github.com/gphotosuploader/google-photos-api-client-go/v3` from 3.0.1 to 3.0.2
- Bump `golang.org/x/oauth2` from 0.12.0 to 0.13.0
- Bump `golang.org/x/sync` from 0.3.0 to 0.4.0 ([#377][i377])
- Bump `golang.org/x/term` from 0.10.0 to 0.13.0 ([#376][i376])
- [CI] Bump `github.com/stretchr/testify` from 1.7.0 to 1.8.4 ([#380][i380])
- [CI] Bump `actions/checkout` from 3 to 4 ([#375][i375])
- [CI] Bump `goreleaser/goreleaser-action` from 4 to 5 ([#374][i374])
- [CI] Bump `golangci` from 1.52.1 to 1.54.2

### Removed
- Support for multiple concurrent workers. The bandwidth to upload items is shared, so we are not expecting any performance problem.
- Removed DEPRECATED configuration parameters from previous versions.

[i374]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/374
[i375]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/375
[i376]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/376
[i377]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/377
[i378]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/378
[i379]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/379
[i380]: https://github.com/gphotosuploader/gphotos-uploader-cli/pulls/380

## 3.5.2
### Added
- Support for the latest published Go version (1.21). This project will maintain compatibility with the latest two major versions published.
- Client cache for albums to reduce the number of requests to Google Photos API


### Changed
- Bump `github.com/sirupsen/logrus` from 1.8.1 to 1.9.3
- Bump `github.com/spf13/afero` from 1.8.2 to 1.10.0 
- Bump `golang.org/x/oauth2` from 0.12.0 to 0.13.0
- Bump `golang.org/x/sync` from 0.3.0 to 0.4.0 
- Bump `golang.org/x/term` from 0.10.0 to 0.13.0
- Bump `github.com/99designs/keyring` from 1.2.1 to 1.2.2
- Bump `github.com/gphotosuploader/google-photos-api-client-go/v2` from 2.4.0 to 2.4.2
- Bump `github.com/schollz/progressbar/v3` from 3.8.6 to 3.13.1
- Bump `github.com/spf13/cobra` from 1.4.0 to 1.7.0
- Bump `golang.org/x/oauth2` from v0.0.0-20220309155454-6242fa91716a to 0.13.0
- Bump `golang.org/x/sync` from v0.0.0-20210220032951-036812b2e83c to 0.4.0
- Bump `golang.org/x/term` v0.0.0-20210927222741-03fcf44c2211 to 0.13.0
- Bump `google.golang.org/api` from v0.74.0 to 0.148.0

## 3.5.1
### Added
- Support for the latest published Go version (1.20). This project will maintain compatibility with the latest two major versions published.

### Fixed
- Restrict allowed SecretsBackendTypes to the ones supported by the CLI. ([#347][i347])

### Removed
- Support for previous Go version (1.18).

[i347]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/347

## 3.5.0
### Added
- Support for the latest published Go version (1.19). This project will maintain compatibility with the latest two major versions published.
### Fixed
- Exit if daily API quota is exceeded.  Thanks to [@mlbright](https://github.com/mlbright) ([#341][i341])
### Removed
- Once Go 1.19 has been published, previous Go 1.17 support is deprecated.

[i341]: https://github.com/gphotosuploader/gphotos-uploader-cli/pull/341

## 3.4.0
### Changed
- The command `auth` initiates the [Google authentication to get an OAuth 2.0 token](https://gphotosuploader.github.io/gphotos-uploader-cli/#/getting-started?id=authentication). **It should be used the first time that the CLI is configured**. See [documentation](https://gphotosuploader.github.io/gphotos-uploader-cli/#/getting-started?id=authentication).
### Deprecated 
- Google deprecates the OAuth 2.0 authentication based on out-of-band tokens. ([#326][i326])

[i326]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/326

## 3.3.1
### Added
- Support for the latest published Go version (1.18). This project will maintain compatibility with the latest two major versions published.
### Changed
- Dependency has been updated, so potential bugs have been fixed.
### Deprecated
- Once Go 1.18 has been published, previous Go 1.16 support is deprecated.
### Removed
- Command `release` in the Makefile. We are using [goreleaser GitHub action](https://github.com/goreleaser/goreleaser-action) now.

## 3.3.0
### Changed
- Files are sorted before being uploaded. This will only be true for the files uploaded to an empty albums. ([#301][i301])

[i301]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/301
 
## 3.2.1
### Added
- Support for the latest published Go version. This project will maintain compatibility with the latest two major versions published.
### Deprecated
- Once Go 1.17 has been published, previous Go 1.15 support is deprecated.
### Fixed
- Using environment var `GPHOTOS_CLI_TOKENSTORE_KEY`, it was not possible to set an empty key. Now, it is.

## 3.2.0
### Changed
- Reduce the cost of tracking already uploaded files by bringing back file last modification time check ([#306][i306])

[i306]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/306

## 3.1.1
### Fixed
- Keychain backend not working on macOS. Thanks to [@mlangenberg](https://github.com/mlangenberg) ([#302][i302])

[i302]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/302

## 3.1.0
### Added
- Removes retry when Google Photos requests quota limit has been reached. ([#290][i290])
- Removes retry when Google Photos requests quota limit has been reached. ([#248][i248])
- Add support for `go v1.16`.
- Bump `golangci-lint` to `v1.39.0`. 
### Fixed
- Not possible to enter a passphrase - panic: crypto/hmac: hash generation function does not produce unique values ([#294][i294])
### Removed
- Remove support for `go v.1.14`. 

[i294]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/294
[i290]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/290
[i248]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/248

## 3.0.1
### Fixed
- Tagged extension matches with uppercase file extensions.  ([#283][i283])

[i283]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/283

## 3.0.0
> This is a **major upgrade**, so it has some **non-backwards compatible changes**
### Added
- **Progress bar** when uploading files.
- Configuration, wo/ sensible data, is printed when debug is enabled. ([#270][i270]) 
- Configuration validation. The cli validates the configuration data at starting time.
- Information messages to bring more context at runtime. ([#260][i260]) 
### Changed
- `Jobs.MakeAlbums` configuration setting has changed to `Jobs.CreateAlbums`. Valid values are `Off`,`folderName` and `folderPath`.
- **Reduce the number of calls to the API when uploading files**. It's using less than 50% of calls than before.
- Move to `golang.org/x/term` from `golang.org/x/crypto/ssh/terminal`, due to deprecation.
- Some parts of the code have been refactored to make cleaner code and increase testability.
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
> - `includePatterns` & `excludePatterns` configuration options have changed.
> - `includePatterns` has a new default (`_IMAGE_EXTENSIONS_`).
> - `uploadVideos` configuration option has been removed.
### Added
- Two new tagged patterns have been added: `_IMAGE_EXTENSIONS_`, matching [supported image file types](https://support.google.com/googleone/answer/6193313), and `_RAW_EXTENSIONS_`, matching [supported RAW file types](https://support.google.com/googleone/answer/6193313). ([#249][i249])
- Retry management. It's implementing exponential back-off with a maximum of 4 retries by default.  ([#253][i253]) 
### Changed
- `includePatterns` & `excludePatterns` configuration options have changed. It's using a new format, please review de [configuration documentation][idocumentation].
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
- Remove `build` from a version. Now `version` has all the tag+build information.
### Fixed
- Fix duplicated album creation. ([#135][i135])

[i135]: https://github.com/gphotosuploader/gphotos-uploader-cli/issues/135

## 0.8.5
### Fixed
- Fix exit code on clean executions. (#137)

## 0.8.4
### Fixed
- Fix panic when an unexpected error on media item creation was raised. (#110)
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
- Fix duplicated album creation due to a concurrency problem. (#19)

## 0.8.0
### Added
- Uploads can be resumed. This will help upload large files or when connection has fails. Thanks to @pdecat.

## 0.7.2
### Fixed
- Fix token storing when expired token has been refreshed. See comments on #107.
- Fix typos in CHANGELOG.

## 0.7.1 
### Added
- Google Auth expired token refresh. Once token is expired, `gphotos-uploader-cli` will try to refresh the token without user intervention. **NOTE**: First time you use this version, you should re-authenticate in order to get the token that allows token refresh. (#103)
- Add `--config` flag to specify the folder where configuration is kept. (#104)
### Changed
- Moved CI/CD platform from Travis to [Drone.io](https://cloud.drone.io/gphotosuploader/gphotos-uploader-cli). It has reduced the time to CI by a half.

## 0.6.0
### Added
- `deleteAfterUpload` option has been reactivated, it was removed on v0.4.0. If you use this option in [config file][idocumentation] files will be deleted from local repository after being uploaded to Google Photos. (#25)
### Changed
- This repository has transferred to [GPhotos Uploaders organization](https://github.com/gphotosuploader), so all imports have been updated to the new organization's URL.
### Removed
- Removed some useless log lines. There is still too much.

## 0.5.0
### Changed
- Fix issue #97 "New gnome keyring store created on each launch". To solve this issue a new `serviceName` has been changed. **NOTE**: Once you use this version, a new Gnome keyring will be created, so credentials should be supplied again. (#97) 

## 0.4.2
### Fixed
- Fix CI release pipeline to fix an application version (#94). The Last version was still broken on CI.

## 0.4.1
### Added
- Add Homebrew tap to allow users to install `gphotos-uploader-cli` using Homebrew. See [install](README.md) section.
 
### Fixed
- Fix CI release pipeline to fix an application version (#94)

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
- `make test` is not as verbose as before. To make it easier to see if there is an error
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
- Fix goreleaser configuration (remove a deprecated statement)
- Update [Getting started](README.md) documentation

### Removed
- Remove [snap](https://snapcraft.io/snaps) application publication (someone has stolen our app name)

## 0.1.11 - 2018-09-20
### Added
- [goreleaser](https://goreleaser.com/) will be in charge of publishing [binaries](https://github.com/gphotosuploader/gphotos-uploader-cli/releases) after the new release is done

[idocumentation]: https://gphotosuploader.github.io/gphotos-uploader-cli/
