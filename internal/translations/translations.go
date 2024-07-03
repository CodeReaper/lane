package translations

import (
	"encoding/csv"
	"os"
	"strings"
)

type translationData struct {
	data [][]string
}

func newTranslations(path string) (*translationData, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(r)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		records = records[1:] // skips header row
	}
	return &translationData{data: records}, nil
}

func (t *translationData) translation(keyIndex int, valueIndex int) *translation {
	items := map[string]string{}
	for _, r := range t.data {
		items[strings.ToLower(r[keyIndex-1])] = r[valueIndex-1]
	}
	return newTranslation(items)
}
