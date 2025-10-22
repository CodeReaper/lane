package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
)

type MockService struct {
	resp  *http.Response
	files *drive.FileList
	err   error
}

func (s MockService) download(documentId string, mimeType string) (*http.Response, error) {
	return s.resp, s.err
}

func (s MockService) list() (*drive.FileList, error) {
	return s.files, s.err
}

var expectedDownloadPath = "../../build/out.csv"

func cleanup() {
	log.Println("setup test")
	os.Remove(expectedDownloadPath)
}

func TestDownloadEmptyCredentialsFailure(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}

	err := Download(context.TODO(), &flags)

	assert.Error(t, err)
}

func TestDownloadNoCredentialsFailure(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/does-not-exist.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}

	err := Download(context.TODO(), &flags)

	assert.Error(t, err)
}

func TestDownloadServiceFailure(t *testing.T) {
	defer cleanup()
	flags := Flags{
		Output:      expectedDownloadPath,
		Credentials: "testdata/empty.json",
		DocumentId:  "1234567890",
		Format:      "csv",
	}
	svc := MockService{
		err: fmt.Errorf("always fails"),
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
	svc := MockService{}
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

func TestListEmptyCredentialsFailure(t *testing.T) {
	_, err := List(context.TODO(), "testdata/empty.json")

	assert.Error(t, err)
}

func TestListNoCredentialsFailure(t *testing.T) {
	_, err := List(context.TODO(), "testdata/does-not-exist.json")

	assert.Error(t, err)
}

func TestListServiceFailure(t *testing.T) {
	svc := MockService{
		err: fmt.Errorf("always fails"),
	}
	r, err := list(svc)
	assert.Nil(t, r)
	assert.Error(t, err)
}

func TestList(t *testing.T) {
	svc := MockService{
		files: &drive.FileList{Files: []*drive.File{{Id: "id"}}},
	}
	r, err := list(svc)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, 1, len(r))
	assert.ElementsMatch(t, []string{"id"}, keys(r))
}
