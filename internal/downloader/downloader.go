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

func Download(ctx context.Context, flags *Flags) error {
	service, err := newService(ctx, flags.Credentials)
	if err != nil {
		return err
	}

	return download(flags, service)
}

func download(flags *Flags, service Service) error {
	if err := flags.validate(); err != nil {
		return err
	}

	mimeType, err := lookupMimeType(flags.Format)
	if err != nil {
		return err
	}

	resp, err := service.download(flags.DocumentId, mimeType)
	if err != nil {
		return err
	}

	return handleResponse(resp, flags.Output)
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
