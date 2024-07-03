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
			Input:             "",
			Kind:              "",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "",
		},
		false,
	},
	{
		"all-set-android",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "android",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "",
		},
		true,
	},
	{
		"all-set-ios",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "testdata/out.put",
		},
		true,
	},
	{
		"all-set-but-unknown-kind",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "unknown",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-input",
		Flags{
			Kind:              "ios",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-kind",
		Flags{
			Input:             "testdata/input.csv",
			KeyIndex:          0,
			DefaultValueIndex: 0,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-output",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          0,
			DefaultValueIndex: 0,
		},
		false,
	},
	{
		"all-set-ios-but-missing-key",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			DefaultValueIndex: 0,
			Output:            "testdata/out.put",
		},
		true,
	},
	{
		"all-set-ios-but-missing-value",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "ios",
			KeyIndex: 0,
			Output:   "testdata/out.put",
		},
		true,
	},
}

func TestFlagsValidate(t *testing.T) {
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
