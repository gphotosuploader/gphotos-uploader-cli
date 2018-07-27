package fileshandling

import "os"

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

// func deleteFile(filePath string)  error {
// 	err := os.Remove(filePath)
// 	if err != nil {
// 		return
// 	}
// 	return false, nil
// }

// func deleteFile() {
// 	// delete file
// 	var err = os.Remove(path)
// 	if isError(err) {
// 		return
// 	}

// 	fmt.Println("==> done deleting file")
// }
