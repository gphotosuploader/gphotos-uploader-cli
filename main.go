package main

import "gitlab.com/nmrshll/gphotos-uploader-go-cookies/config"

const (
	imagePath = "/home/me/photos_autres/USSIS/2014_11_WE_U6/DSC_0501.JPG"
)

func main() {
	// load all config parameters
	uploadJobs := config.Load()

	for _, job := range uploadJobs {
		job.Run()
	}
}
