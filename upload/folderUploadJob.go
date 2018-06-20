package upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/nmrshll/gphotos-uploader-go-cookies/filesHandling"

	"github.com/simonedegiacomi/gphotosuploader/auth"
)

const (
	USEFOLDERNAMES = "folderNames"
)

type FolderUploadJob struct {
	Credentials  *auth.CookieCredentials
	Account      string
	SourceFolder string
	MakeAlbums   struct {
		Enabled bool
		Use     string
	}
	DeleteAfterUpload bool
}

func (j *FolderUploadJob) Run() {
	err := j.uploadFolder(j.SourceFolder)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (j *FolderUploadJob) uploadFolder(folderPath string) error {
	if !filesHandling.IsDir(folderPath) {
		return fmt.Errorf("%s is not a folder", folderPath)
	}

	defer close(fileUploadsChan)
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filesHandling.IsFile(path) {
			var fu = &FileUpload{FolderUploadJob: j, filePath: path}
			if j.MakeAlbums.Enabled && j.MakeAlbums.Use == USEFOLDERNAMES {
				lastDirName := filepath.Base(filepath.Dir(path))
				fu.albumName = lastDirName
			}
			fileUploadsChan <- fu
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
