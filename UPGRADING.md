# Upgrading notes

## Upgrading To 2.0.0 from 1.x.x

### Patterns definition

The `includePatterns` and `excludePatterns` configuration options has changed, see [configuration documentation](.docs/configuration.md). You should modify your configuration to honor the **new format**.

If you were using the tagged patterns (*\_ALL_FILES_* and *\_ALL_VIDEO_FILES_*) you don't need to do anything. 

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
 