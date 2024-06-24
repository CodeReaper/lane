package translations

import (
	"encoding/csv"
	"os"
	"strings"
)

type Translations struct {
	data [][]string
}

func newTranslations(path string) (*Translations, error) {
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
	return &Translations{data: records}, nil
}

func (t *Translations) translation(keyIndex int, valueIndex int) *Translation {
	items := map[string]string{}
	for _, r := range t.data {
		items[strings.ToLower(r[keyIndex-1])] = r[valueIndex-1]
	}
	return newTranslation(items)
}
