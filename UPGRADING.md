# Upgrading notes

## Upgrading To 3.x from 2.x

### Configuration settings
- `Jobs.Account` configuration setting has been changed to `Account`. See [configuration documentation](https://gphotosuploader.github.io/gphotos-uploader-cli/#/configuration).
- `Jobs.MakeAlbums` configuration setting has changed to `Jobs.CreateAlbums`. See [configuration documentation](https://gphotosuploader.github.io/gphotos-uploader-cli/#/configuration?id=createalbums).
- **Multiple Google Photos account support has been removed**. You can use multiple configuration files in the same application folder instead.

## Upgrading To 2.x from 1.x

### Patterns definition

The `includePatterns` and `excludePatterns` configuration options has changed, see [configuration documentation](https://gphotosuploader.github.io/gphotos-uploader-cli/#/configuration). You should modify your configuration to honor the **new format**.

If you were using the tagged patterns (`_ALL_FILES_` and `_ALL_VIDEO_FILES_`) you don't need to do anything. 

```bash
sourceFolder
`-- foo
    |-- picture1.png
    |-- picture2.png
    `-- bar
        |-- picture1.png
        |-- picture2.png
```
#### Some examples
Description | Current format | Previous format
----------- | -------------- | ---------------
Include all files | `includePatterns: "**"}` | `includePatterns: {"*"}`
Include only PNG files | `includePatterns: "**/*.png"}` | `includePatterns: {"*.png"}`
Include PNG files in `foo` folder | `includePatterns: "foo/*.png"}` | `includePatterns: {"*.png"}` <br> `excludePatterns: {"bar"}`
 
