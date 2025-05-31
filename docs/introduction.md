# Introduction

`gphotos-uploader-cli` is a command-line tool for mass uploading media folders to your Google Photos account.

While the official Google Photos uploader only supports macOS and Windows, this tool brings uploading capabilities to Linux and any OS where you can compile a Go program.

Some notes before you start:

* This tool is not affiliated with Google. It is an independent project that uses the [Google Photos API](https://developers.google.com/photos).
* The Google Photos API has several limitations which this tool can solve. Please read the [Limitations](#limitations) section to ensure it meets your needs.
* By the nature of how this CLI has been designed, **it is not suitable to be used inside a Docker container**. It is designed to run on a local machine with a persistent filesystem and access to the interactive terminal for authentication.

## Features

- **Customizable configuration:** Use a JSON-like config file.
- **Flexible file filtering:** Include or exclude files and folders using patterns (see [configuration documentation](configuration.md)).
- **Resumable uploads:** Resume interrupted uploads to save time and bandwidth.
- **Automatic file deletion:** Optionally delete local files after uploading.
- **Smart upload tracking:** Only new files are uploaded, saving bandwidth.
- **Local caching:** Reduces the number of queries to Google Photos.
- **Secure authentication:** Uses OAuth for login; stores access tokens in your OS's secure storage (keyring/keychain).
- **Robust retry logic:** All requests are retried with exponential back-off, following [Google Photos best practices](https://developers.google.com/photos/library/guides/best-practices#error-handling).

## Limitations

- **Supported file types:** Only images and videos can be uploaded. Unsupported formats will be uploaded, but Google Photos will reject them as media items.
- **Photo storage and quality:** All uploads are stored in full, [original quality](https://support.google.com/photos/answer/6220791) and count toward your storage quota. "High quality" mode is not available via the API.
- **Duplicates:** Uploading the same file twice results in deduplication, but the original filename is retained.
- **Media dates:** The date shown in Google Photos is based on EXIF creation date or upload date, not the local file's modification date. This cannot be changed by `gphotos-uploader-cli`.
- **File size checks:** The API does not return media sizes. Existence checks are possible, but size checks require slow, additional HTTP requests.
- **Album management:** Files can only be uploaded to albums created by `gphotos-uploader-cli`. Only files uploaded by this tool can be removed from those albums.
- **API rate limits:** Google Photos enforces daily quotas:
  - 10,000 requests per project per day for the Library API
  - 75,000 requests per project per day for media bytes access
- **Not suitable for Docker:** This CLI is designed for local execution with persistent storage and interactive terminal access, making it unsuitable for containerized environments.