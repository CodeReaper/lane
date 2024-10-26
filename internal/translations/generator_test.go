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

func TestAndroidWithFillIn(t *testing.T) {
	flags := Flags{
		Input:             "testdata/fill-in/input.csv",
		Kind:              "android",
		KeyIndex:          1,
		DefaultValueIndex: 2,
		FillIn:            true,
	}
	configurations := []string{"2 ../../build/en.xml", "3 ../../build/da.xml"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/fill-in/android-en.expected", "../../build/en.xml")
	equalFiles(t, "testdata/fill-in/android-da.expected", "../../build/da.xml")
}

func TestJsonWithFillIn(t *testing.T) {
	flags := Flags{
		Input:             "testdata/fill-in/input.csv",
		Kind:              "json",
		KeyIndex:          1,
		DefaultValueIndex: 2,
		FillIn:            true,
	}
	configurations := []string{"2 ../../build/en.json", "3 ../../build/da.json"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/fill-in/json-en.expected", "../../build/en.json")
	equalFiles(t, "testdata/fill-in/json-da.expected", "../../build/da.json")
}

func TestIosWithFillIn(t *testing.T) {
	flags := Flags{
		Input:             "testdata/fill-in/input.csv",
		Kind:              "ios",
		KeyIndex:          1,
		DefaultValueIndex: 2,
		Output:            "../../build/Translations.swift",
		FillIn:            true,
	}
	configurations := []string{"2 ../../build/en.strings", "3 ../../build/da.strings"}

	err := Generate(context.Background(), &flags, configurations)

	assert.Nil(t, err)
	equalFiles(t, "testdata/fill-in/ios-en.expected", "../../build/en.strings")
	equalFiles(t, "testdata/fill-in/ios-da.expected", "../../build/da.strings")
	equalFiles(t, "testdata/fill-in/ios-swift.expected", "../../build/Translations.swift")
}

func TestInvalidConfigurationIndex(t *testing.T) {
	flags := Flags{
		Input:    "testdata/input.csv",
		Kind:     "json",
		KeyIndex: 1,
	}
	configurations := []string{"0 ../../build/en.json"}

	err := Generate(context.Background(), &flags, configurations)

	assert.NotNil(t, err)
}

var configurationCases = []struct {
	name           string
	configurations []string
	passes         bool
}{
	{
		"none-set",
		make([]string, 0),
		false,
	},
	{
		"invalid-index",
		[]string{"0 ../../build/en.json"},
		false,
	},
	{
		"missing-index",
		[]string{"0 ../../build/en.json"},
		false,
	},
	{
		"missing-path",
		[]string{"0"},
		false,
	},
	{
		"all-good",
		[]string{"1 ../../build/en.json"},
		true,
	},
}

func TestConfigurationCases(t *testing.T) {
	flags := Flags{
		Input:    "testdata/input.csv",
		Kind:     "json",
		KeyIndex: 1,
	}

	for _, c := range configurationCases {
		t.Run(c.name, func(t *testing.T) {
			err := Generate(context.Background(), &flags, c.configurations)
			if c.passes && err != nil {
				t.Errorf("expected to pass, but got %v", err)
			}
			if !c.passes && err == nil {
				t.Errorf("expected to fail, but got %v", err)
			}
		})
	}
}
