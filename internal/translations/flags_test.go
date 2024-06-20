package translations

import (
	"testing"
)

var validationKases = []struct {
	name   string
	flags  Flags
	passes bool
}{
	{
		"none-set",
		Flags{},
		false,
	},
	{
		"all-empty",
		Flags{
			Input:          "",
			Kind:           "",
			Index:          0,
			Configurations: []string{},
			DefaultIndex:   0,
			Output:         "",
		},
		false,
	},
	{
		"all-set-android",
		Flags{
			Input:          "testdata/input.csv",
			Kind:           "android",
			Index:          0,
			Configurations: []string{""},
			DefaultIndex:   0,
			Output:         "",
		},
		true,
	},
	{
		"all-set-ios",
		Flags{
			Input:          "testdata/input.csv",
			Kind:           "ios",
			Index:          0,
			Configurations: []string{""},
			DefaultIndex:   0,
			Output:         "",
		},
		true,
	},
	// FIXME: more test cases, and fix confs
}

func TestConfigValidate(t *testing.T) {
	for _, kase := range validationKases {
		t.Run(kase.name, func(t *testing.T) {
			err := kase.flags.validate()
			if kase.passes && err != nil {
				t.Errorf("expected to pass, but got %v", err)
			}
			if !kase.passes && err == nil {
				t.Errorf("expected to fail, but got %v", err)
			}
		})
	}
}
