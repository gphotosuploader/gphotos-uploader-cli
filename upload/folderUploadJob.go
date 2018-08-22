package upload

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"

	"github.com/nmrshll/go-cp"
	gphotos "github.com/nmrshll/google-photos-api-client-go/noserver-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
)

const (
	USEFOLDERNAMES = "folderNames"
)

type FolderUploadJob struct {
	*config.FolderUploadJob
}

func (folderUploadJob *FolderUploadJob) Run() {
	sourceFolderAbsolutePath, err := cp.AbsolutePath(folderUploadJob.SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	client, err := Authenticate(folderUploadJob)
	if err != nil {
		log.Fatal(err)
	}

	err = folderUploadJob.uploadFolder(client, sourceFolderAbsolutePath)
	if err != nil {
		log.Fatal(err)
	}
}

func Authenticate(folderUploadJob *FolderUploadJob) (*gphotos.Client, error) {
	var httpClient *http.Client

	// try to load token from keyring
	token, err := tokenstore.RetrieveToken(folderUploadJob.Account)
	if err == nil {
		// if found create client from token
		// httpClient = gphotosapiclient.NewClientFromToken(token)
		gphotosClient, err := gphotos.NewClient(gphotos.FromToken(config.OAuthConfig(), token))
		if err == nil && gphotosClient != nil {
			return gphotosClient, nil
		}
	} else {
		// else whatever the reason authenticate again to grab a new token
		authorizedClient, err := gphotosapiclient.NewOAuthClient()
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed authenticating new client")
		}

		// and store the token into the keyring
		err = tokenstore.StoreToken(folderUploadJob.Account, *authorizedClient.Token)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed storing token")
		}

		httpClient = authorizedClient.Client
	}
	if httpClient == nil {
		return nil, stacktrace.NewError("httpClient shouldn't be still nil")
	}

	photosClient, err := gphotosapiclient.New(httpClient)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed creating new photos client from httpClient")
	}
	return photosClient, nil
}

func (j *FolderUploadJob) uploadFolder(gphotosClient *gphotosapiclient.PhotosClient, folderPath string) error {
	if !fileshandling.IsDir(folderPath) {
		return fmt.Errorf("%s is not a folder", folderPath)
	}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if fileshandling.IsFile(path) {
			var fileUpload = &FileUpload{FolderUploadJob: j, filePath: path, gphotosClient: *gphotosClient}
			if j.MakeAlbums.Enabled && j.MakeAlbums.Use == USEFOLDERNAMES {
				lastDirName := filepath.Base(filepath.Dir(path))
				fileUpload.albumName = lastDirName
			}
			QueueFileUpload(fileUpload)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return nil
}

// const reduceGetFirst = (a, v) => (!a && v ? v : a);

// const uploadFile = async(gphotos, filePath, albumName, context) => {
//     const uploadedImage = await gphotos.upload(filePath).catch(logError);
//     console.log(filePath, albumName, context);
//     if (context.makeAlbums && albumName) {
//         const album = await gphotos.searchOrCreateAlbum(albumName);
//         await album.addPhoto(uploadedImage);
//     }
//     if (context.deleteAfterUpload)
//         checkAndDeleteLocal(uploadedImage.rawUrl, filePath).catch(
//             wrapLogError("checkAndDeleteLocal")
//         );
// };

// const existsPath = path =>
//     fs.access(path, function(err) {
//         if (err && err.code === "ENOENT") {
//             console.log("this runs");
//             return false;
//         }
//         return true;
//     });

// function isDirSync(aPath) {
//     try {
//         return fs.statSync(aPath).isDirectory();
//     } catch (e) {
//         if (e.code === "ENOENT") {
//             return false;
//         } else {
//             throw e;
//         }
//     }
// }

// const uploadFolder = (gphotos, folderPath, context) => {
//     if (isDirSync(folderPath)) {
//         const walker = walk.walk(folderPath);

//         walker.on("file", (root, fileStats, next) => {
//             fs.readFile(fileStats.name, () => {
//                 const filePath = `${root}/${fileStats.name}`;
//                 const firstSubFolder = root
//                     .replace(folderPath, "")
//                     .split("/")
//                     .reduce(reduceGetFirst);
//                 if (!context.exclude.endsWith.some(v => filePath.endsWith(v))) {
//                     uploadFile(gphotos, filePath, (albumName = firstSubFolder), context)
//                         .then(uploadedImage => {
//                             next();
//                         })
//                         .catch(wrapLogError("uploadFile"));
//                 } else {
//                     console.log(`skipping file ${filePath}: excluded file extension`);
//                 }
//             });
//         });

//         walker.on("errors", function(root, nodeStatsArray, next) {
//             nodeStatsArray.map(stat => wrapLogError("walker")(stat.error));
//             next();
//         });

//         walker.on("end", function() {
//             console.log("all done");
//         });
//     } else {
//         console.log("path doesn't exist", folderPath);
//     }
// };
