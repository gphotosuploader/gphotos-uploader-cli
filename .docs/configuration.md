## Configuration

Example configuration file:
[embedmd]:# (../config/config.example.hjson)
```hjson
{
  APIAppCredentials: {
    ClientID:     "20637643488-1hvg8ev08r4tc16ca7j9oj3686lcf0el.apps.googleusercontent.com",
    ClientSecret: "0JyfLYw0kyDcJO-pGg5-rW_P",
  }
  jobs: [
    {
      account: youremail@gmail.com
      sourceFolder: ~/folder/to/upload
      makeAlbums: {
        enabled: true
        use: folderNames
      }
      deleteAfterUpload: true
    }
  ]
}
```

# `APIAppCredentials`:
The credentials that are provided are just example ones. 
Replace them with credentials you create at https://console.cloud.google.com/apis/api/photoslibrary.googleapis.com

# `jobs`:
List of folders to upload and upload options for each folder.

## `account`:
Needs to be unique.
If it contains a google email address, it will be suggested at login.

