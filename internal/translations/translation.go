package translations

import "slices"

type Translation struct {
	keys  []string
	items map[string]string
}

func newTranslation(items map[string]string) *Translation {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	return &Translation{
		keys:  keys,
		items: items,
	}
}

func (t *Translation) get(key string) string {
	return t.items[key]
}
