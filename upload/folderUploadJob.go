package upload

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"

	"github.com/davecgh/go-spew/spew"
	"github.com/nmrshll/go-cp"

	"gitlab.com/nmrshll/gphotos-uploader-go-api/config"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/datastore"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/filesHandling"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/gphotosapiclient"
)

const (
	USEFOLDERNAMES = "folderNames"
)

// var (
// 	folderUploadsChan = make(chan struct{})
// )

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

func Authenticate(folderUploadJob *FolderUploadJob) (*gphotosapiclient.PhotosClient, error) {
	var httpClient *http.Client

	// try to load token from keyring
	token, err := tokenstore.RetrieveToken(folderUploadJob.Account)
	if err == nil {
		// if found create client from token
		httpClient = gphotosapiclient.NewClientFromToken(token)
	} else {
		// else whatever the reason authenticate again to grab a new token
		authorizedClient, err := gphotosapiclient.NewOAuthClient()
		if err != nil {
			log.Fatal(err)
		}

		// and store the token into the keyring
		err = tokenstore.StoreToken(folderUploadJob.Account, *authorizedClient.Token)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed storing token")
		}

		httpClient = authorizedClient.Client
	}
	// if err != nil && err != tokenstore.ErrNotFound {
	// 	fmt.Println("failed retrieving valid token, browser will open to authenticate again...")
	// }

	// // if found create a client from the token
	// if err == nil {
	// 	// create client from token
	// }

	// // if not found authenticate the user to create a new OAuthClient and store token in keyring
	// if err == tokenstore.ErrNotFound {
	// }

	photosClient, err := gphotosapiclient.New(httpClient)
	if err != nil {
		log.Fatal(err)
	}
	return photosClient, nil
}

func (j *FolderUploadJob) uploadFolder(gphotosClient *gphotosapiclient.PhotosClient, folderPath string) error {
	if !filesHandling.IsDir(folderPath) {
		return fmt.Errorf("%s is not a folder", folderPath)
	}

	// defer close(fileUploadsChan)
	// NO. TODO: USE ONE CHAN PER FOLDERUPLOADJOB

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		spew.Dump(path)
		if info.IsDir() {
			return nil
		}
		if filesHandling.IsFile(path) {
			var fileUpload = &FileUpload{FolderUploadJob: j, filePath: path, gphotosClient: *gphotosClient}
			if j.MakeAlbums.Enabled && j.MakeAlbums.Use == USEFOLDERNAMES {
				lastDirName := filepath.Base(filepath.Dir(path))
				fileUpload.albumName = lastDirName
			}
			QueueFileUpload(fileUpload)
		}
		// if filepath.Ext(path) == ".sh" {
		// 	list = append(list, path)
		// }

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
