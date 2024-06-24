package translations

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type LanguageFile struct {
	path       string
	keyIndex   int
	valueIndex int
}

func newLanguageFiles(keyIndex int, configurations []string) ([]LanguageFile, error) {
	if len(configurations) == 0 {
		return nil, fmt.Errorf("no configurations provided")
	}

	list := make([]LanguageFile, 0)
	for _, configuration := range configurations {
		fields := strings.Fields(configuration)
		if len(fields) != 2 {
			return nil, fmt.Errorf("configuration has invalid format: %s", configuration)
		}

		index, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("configuration has invalid index: %s", configuration)
		}

		path := fields[1]
		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			return nil, fmt.Errorf("configuration has invalid path: %s", configuration)
		}

		list = append(list, LanguageFile{
			path:       path,
			keyIndex:   keyIndex,
			valueIndex: index,
		})
	}
	return list, nil
}
