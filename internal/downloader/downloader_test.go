package downloader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockService struct {
	resp *http.Response
	err  error
}

func (s MockService) download(documentId string, mimeType string) (*http.Response, error) {
	return s.resp, s.err
}

var expectedDownloadPath = "../../build/out.csv"

func cleanup() {
	log.Println("setup test")
	os.Remove(expectedDownloadPath)
}

func TestDownloadFailure(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}
	svc := MockService{
		resp: nil,
		err:  fmt.Errorf("always fails"),
	}
	err := download(&flags, svc)
	assert.Error(t, err)
}

func TestDownloadFailedResponse(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}
	svc := MockService{
		resp: &http.Response{
			Status:        "Forbidden",
			StatusCode:    401,
			Body:          io.NopCloser(&bytes.Buffer{}),
			ContentLength: 0,
		},
		err: nil,
	}
	err := download(&flags, svc)
	assert.Error(t, err)
	assert.NoFileExists(t, expectedDownloadPath)
}

func TestDownloadInvalidFlags(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "unknown",
	}
	svc := MockService{
		resp: nil,
		err:  nil,
	}
	err := download(&flags, svc)
	assert.Error(t, err)
}

func TestDownload(t *testing.T) {
	defer cleanup()
	var expected = []byte("test,data")

	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}
	svc := MockService{
		resp: &http.Response{
			Status:        "OK",
			StatusCode:    200,
			Body:          io.NopCloser(bytes.NewBuffer(expected)),
			ContentLength: int64(len(expected)),
		},
		err: nil,
	}
	err := download(&flags, svc)
	assert.NoError(t, err)

	b, err := os.ReadFile(expectedDownloadPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)
}

func TestMimetypeLookup(t *testing.T) {
	m, err := lookupMimeType("csv")
	assert.NotEmpty(t, m)
	assert.NoError(t, err)

	m, err = lookupMimeType("unknown")
	assert.Empty(t, m)
	assert.Error(t, err)
}
