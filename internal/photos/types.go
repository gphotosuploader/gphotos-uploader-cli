package photos

import (
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/oauth2/google"
)

// Endpoint is an URL of Google Photos Library API.
var Endpoint = google.Endpoint

// Scopes is a set of OAuth scopes.
var Scopes = []string{photoslibrary.PhotoslibraryScope}
