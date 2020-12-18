# Introduction

Command line tool to mass upload media folders to your Google Photos account.    

While the official tool only supports Mac OS and Windows, this brings an uploader to Linux too. Lets you upload photos from, in theory, any OS for which you can compile a Go program.    
> The Google Photos API which `gphotos-uploader-cli` uses has quite a few limitations, so please read the [limitations section](#limitations) carefully to make sure it is suitable for your use. 

## Features

- **Customizable configuration**: via JSON-like config file.
- **Filter files with patterns**: include/exclude files & folders using patterns (see [documentation](configuration.md)).
- **Resumable uploads**: Uploads can be resumed, saving time and bandwidth. 
- **File deletion after uploading**: Clean up local files after being uploaded.
- **Track already uploaded files**: uploads only new files to save bandwidth.
- **Caches request results**: keep a local cache to reduce number of queries to Google Photos.
- **Secure**: logs you into Google using OAuth (so this app doesn't have to know your password), and stores your temporary access code in your OS's secure storage (keyring/keychain).
- **Retryable**: all the requests are retried using exponential back-off as is recommended by [Google Photos best practices](https://developers.google.com/photos/library/guides/best-practices#error-handling).

## Limitations
Only images and videos can be uploaded. If you attempt to upload non videos or images or formats that Google Photos doesn't understand, `gphotos-uploader-cli` will upload the file, then Google Photos will give an error when it is put turned into a media item.

### Photo storage and quality
All media items uploaded to Google Photos using the API [are stored in full resolution](https://support.google.com/photos/answer/6220791) at original quality. **They count toward the user’s storage**. The API does not offer a way to upload in "high quality" mode.

### Duplicates
If you upload the same image (with the same binary data) twice then Google Photos will deduplicate it. However it will retain the filename from the first upload which may be confusing. In practise this shouldn't cause too many problems.

### Modified time
The date shown of media in Google Photos is the creation date as determined by the EXIF information, or the upload date if that is not known.
This is not changeable by `gphotos-upload-cli` and is not the modification date of the media on local disk. This means that this CLI cannot use the dates from Google Photos for syncing purposes.

### Size
The Google Photos API does not return the size of media. This means that when syncing to Google Photos, `gphotos-uploader-cli` can only do a file existence check.
It is possible to read the size of the media, but this needs an extra HTTP HEAD request per media item so is **very slow** and uses up a lot of transactions.

### Albums
`gphotos-uploader-cli` can only upload files to albums it created. This is a limitation of the Google Photos API.

`gphotos-uploader-cli` can remove files it uploaded from albums it created only.

### Rate Limiting
Google Photos imposes a rate limit on all API clients. The quota limit for requests to the Library API is 10,000 requests per project per day. The quota limit for requests to access media bytes (by loading a photo or video from a base URL) is 75,000 requests per project per day.
