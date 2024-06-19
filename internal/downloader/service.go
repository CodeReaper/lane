package downloader

import (
	"context"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Service interface {
	download(documentId string, mimeType string) (*http.Response, error)
}

type GoogleAPIService struct {
	ctx     context.Context
	service *drive.Service
}

func newService(ctx context.Context, bytes []byte) (Service, error) {
	config, err := google.JWTConfigFromJSON(bytes, drive.DriveReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)

	drive, err := drive.NewService(ctx, option.WithHTTPClient(client))
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
