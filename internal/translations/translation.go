package translations

import "slices"

type translation struct {
	keys  []string
	items map[string]string
}

func newTranslation(items map[string]string) *translation {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	return &translation{
		keys:  keys,
		items: items,
	}
}

func (t *translation) get(key string) string {
	return t.items[key]
}
