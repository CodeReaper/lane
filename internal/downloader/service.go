package downloader

import (
	"context"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Service interface {
	download(documentId string, mimeType string) (*http.Response, error)
	list() (*drive.FileList, error)
}

type GoogleAPIService struct {
	ctx     context.Context
	service *drive.Service
}

func newService(ctx context.Context, credentials string) (Service, error) {
	drive, err := drive.NewService(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		return nil, err
	}

	return &GoogleAPIService{
		ctx:     ctx,
		service: drive,
	}, nil
}

func (s *GoogleAPIService) download(documentId string, mimeType string) (*http.Response, error) {
	return s.service.Files.Export(documentId, mimeType).Context(s.ctx).Download()
}

func (s *GoogleAPIService) list() (*drive.FileList, error) {
	return s.service.Files.List().Context(s.ctx).Do()
}
