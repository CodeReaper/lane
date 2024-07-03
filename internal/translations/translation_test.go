package translations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTranslation(t *testing.T) {
	items := map[string]string{
		"c": "",
		"d": "",
		"f": "",
		"e": "",
		"a": "",
		"b": "",
	}
	keys := []string{
		"a", "b", "c", "d", "e", "f",
	}

	translation := newTranslation(items)

	assert.Equal(t, items, translation.items)
	assert.EqualValues(t, keys, translation.keys)
}

func TestTranslationGet(t *testing.T) {
	key := "d"
	expected := "4"
	items := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		key: expected,
		"e": "5",
		"f": "6",
	}

	translation := newTranslation(items)

	assert.Equal(t, expected, translation.get(key))
}
