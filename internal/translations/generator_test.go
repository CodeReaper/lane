package translations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAndroid(t *testing.T) {
	flags := Flags{
		Input:    "testdata/input.csv",
		Kind:     "android",
		KeyIndex: 1,
	}
	configurations := []string{"3 ../../build/en.xml", "4 ../../build/da.xml"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/android-en.expected", "../../build/en.xml")
	equalFiles(t, "testdata/android-da.expected", "../../build/da.xml")
}

func TestJson(t *testing.T) {
	flags := Flags{
		Input:    "testdata/input.csv",
		Kind:     "json",
		KeyIndex: 1,
	}
	configurations := []string{"3 ../../build/en.json"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/json-en.expected", "../../build/en.json")
}

func TestIos(t *testing.T) {
	flags := Flags{
		Input:             "testdata/input.csv",
		Kind:              "ios",
		KeyIndex:          1,
		DefaultValueIndex: 3,
		Output:            "../../build/Translations.swift",
	}
	configurations := []string{"3 ../../build/en.strings", "4 ../../build/da.strings"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/ios-en.expected", "../../build/en.strings")
	equalFiles(t, "testdata/ios-da.expected", "../../build/da.strings")
	equalFiles(t, "testdata/ios-swift.expected", "../../build/Translations.swift")
}
