package mock

import (
	"context"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

type GPhotosClient struct {
	ListAlbumsFn      func(ctx context.Context) ([]*photoslibrary.Album, error)
	ListAlbumsInvoked bool

	ListAlbumsWithCallbackFn      func(ctx context.Context, callback gphotos.ListAlbumsFunc) error
	ListAlbumsWithCallbackInvoked bool

	CreateAlbumFn      func(ctx context.Context, title string) (*photoslibrary.Album, error)
	CreateAlbumInvoked bool

	FindAlbumFn      func(ctx context.Context, title string) (*photoslibrary.Album, error)
	FindAlbumInvoked bool

	AddMediaToLibraryFn      func(ctx context.Context, item gphotos.UploadItem) (*photoslibrary.MediaItem, error)
	AddMediaToLibraryInvoked bool

	AddMediaToAlbumFn      func(ctx context.Context, item gphotos.UploadItem, album *photoslibrary.Album) (*photoslibrary.MediaItem, error)
	AddMediaToAlbumInvoked bool
}

// ListAlbums invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) ListAlbums(ctx context.Context) ([]*photoslibrary.Album, error) {
	g.ListAlbumsInvoked = true
	return g.ListAlbumsFn(ctx)
}

// ListAlbumsWithCallback invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) ListAlbumsWithCallback(ctx context.Context, callback gphotos.ListAlbumsFunc) error {
	g.ListAlbumsWithCallbackInvoked = true
	return g.ListAlbumsWithCallbackFn(ctx, callback)
}

// CreateAlbum invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) CreateAlbum(ctx context.Context, title string) (*photoslibrary.Album, error) {
	g.CreateAlbumInvoked = true
	return g.CreateAlbumFn(ctx, title)
}

// FindAlbum invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) FindAlbum(ctx context.Context, title string) (*photoslibrary.Album, error) {
	g.FindAlbumInvoked = true
	return g.FindAlbumFn(ctx, title)
}

// AddMediaToLibrary invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) AddMediaToLibrary(ctx context.Context, item gphotos.UploadItem) (*photoslibrary.MediaItem, error) {
	g.AddMediaToLibraryInvoked = true
	return g.AddMediaToLibraryFn(ctx, item)
}

// AddMediaToAlbum invokes the mock implementation and marks the function as invoked.
func (g *GPhotosClient) AddMediaToAlbum(ctx context.Context, item gphotos.UploadItem, album *photoslibrary.Album) (*photoslibrary.MediaItem, error) {
	g.AddMediaToAlbumInvoked = true
	return g.AddMediaToAlbumFn(ctx, item, album)
}
