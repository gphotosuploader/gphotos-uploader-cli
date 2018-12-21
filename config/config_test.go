package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigFile(t *testing.T) {
	data := []byte(`{
APIAppCredentials: 
{
	ClientID:     "my-client-id",
    ClientSecret: "my-client-secret",
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
`)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	if err != nil {
		t.Errorf("could not create test config file (path: %s): %v", path, err)
	}
	defer func() {
		err := fh.Close()
		if err != nil {
			t.Errorf("could not close test config file: %v", err)
		}
		err = os.Remove(path)
		if err != nil {
			t.Errorf("could not remove test config file (path: %s): %v", path, err)
		}
	}()

	_, err = fh.Write(data)
	if err != nil {
		t.Errorf("could not write data to test config file: %v", err)
	}



}
