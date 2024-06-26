package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Flags struct {
	Output      string
	Credentials string
	DocumentId  string
	Format      string
}

func (f *Flags) validate() error {
	if len(f.Format) == 0 {
		f.Format = "csv"
	}
	if len(f.Output) == 0 {
		return fmt.Errorf("output not provided")
	}
	if len(f.Credentials) == 0 {
		return fmt.Errorf("key not provided")
	}
	if len(f.DocumentId) == 0 {
		return fmt.Errorf("document id not provided")
	}

	validFormat := false
	keys := keys(validFormats)
	for _, v := range keys {
		if !validFormat && v == strings.ToLower(f.Format) {
			validFormat = true
		}
	}
	if !validFormat {
		return fmt.Errorf("invalid format: %s. Valid formats are %v", f.Format, keys)
	}

	if _, err := os.Stat(f.Credentials); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Dir(f.Output)); err != nil {
		return err
	}

	return nil
}
