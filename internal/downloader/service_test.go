package downloader

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func TestNewGoogleAPI(t *testing.T) {
	svc, err := newService(context.Background(), "testdata/empty.json")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

func TestGoogleAPIDownload(t *testing.T) {
	signal := []byte("all good")
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(signal)
	}))
	defer ts.Close()
	svc, err := drive.NewService(ctx, option.WithoutAuthentication(), option.WithEndpoint(ts.URL))
	assert.NoError(t, err, "unable to create client")

	gas := GoogleAPIService{
		ctx:     ctx,
		service: svc,
	}

	resp, err := gas.download("some-id", "a/mime/type")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.EqualValues(t, signal, bytes)
}
