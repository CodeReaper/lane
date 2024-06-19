package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var validFormats = map[string]string{
	"csv":  "text/csv",
	"tsv":  "text/tab-separated-values",
	"ods":  "application/x-vnd.oasis.opendocument.spreadsheet",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

type Client struct {
	flags *Flags
}

func NewClient(x *Flags) *Client {
	return &Client{
		flags: x,
	}
}

func (c *Client) Download(ctx context.Context) error {
	if err := c.flags.validate(); err != nil {
		return err
	}

	mimeType, err := lookupMimeType(c.flags.Format)
	if err != nil {
		return err
	}

	keyBytes, err := os.ReadFile(c.flags.Credentials)
	if err != nil {
		return err
	}

	service, err := newService(ctx, keyBytes)
	if err != nil {
		return err
	}

	resp, err := service.download(c.flags.DocumentId, mimeType)
	if err != nil {
		return err
	}

	return handleResponse(resp, c.flags.Output)
}

func lookupMimeType(format string) (string, error) {
	mimeType, ok := validFormats[strings.ToLower(format)]
	if !ok {
		return "", fmt.Errorf("invalid format: %s. Valid formats are %v", format, keys(validFormats))
	}
	return mimeType, nil
}

func handleResponse(resp *http.Response, outputPath string) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file with http status %d", resp.StatusCode)
	}

	tempPath := outputPath + ".tmp"
	defer os.Remove(tempPath)

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename(tempPath, outputPath)
	if err != nil {
		return err
	}
	return nil
}

func keys(m map[string]string) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
