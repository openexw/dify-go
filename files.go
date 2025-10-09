package dify

import (
	"context"
	"mime/multipart"

	"github.com/go-resty/resty/v2"
	filesv1 "github.com/openexw/dify-go/api/files"
)

type Files interface {
	// UploadFile uploads a file to the server.
	UploadFile(ctx context.Context, req *filesv1.UploadRequest) (*filesv1.UploadResponse, error)
	// PreviewFile previews a file.
	PreviewFile(ctx context.Context, req *filesv1.PreviewRequest) (*filesv1.PreviewResponse, error)
}

type files struct {
	rest   *resty.Client
	appKey string
	multipart.File
}

func NewFiles(rest *resty.Client, appKey string) Files {
	return &files{rest: rest, appKey: appKey}
}

func (f *files) UploadFile(ctx context.Context, req *filesv1.UploadRequest) (*filesv1.UploadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f *files) PreviewFile(ctx context.Context, req *filesv1.PreviewRequest) (*filesv1.PreviewResponse, error) {
	//TODO implement me
	panic("implement me")
}
