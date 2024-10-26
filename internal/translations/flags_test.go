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
			Input:    "testdata/input.csv",
			Kind:     "android",
			KeyIndex: 1,
			Output:   "",
		},
		true,
	},
	{
		"all-set-ios",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
		},
		true,
	},
	{
		"all-set-but-unknown-kind",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "unknown",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-input",
		Flags{
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-kind",
		Flags{
			Input:             "testdata/input.csv",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-output",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
		},
		false,
	},
	{
		"all-set-ios-but-missing-key",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
		},
		false,
	},
	{
		"all-set-ios-but-missing-value",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "ios",
			KeyIndex: 1,
			Output:   "testdata/out.put",
		},
		false,
	},
	{
		"android-with-fill-in",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "android",
			KeyIndex: 1,
			Output:   "",
			FillIn:   true,
		},
		false,
	},
	{
		"android-with-fill-in-and-main-index",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "android",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "",
			FillIn:            true,
		},
		true,
	},
	{
		"ios-with-fill-in-without-main-index",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "ios",
			KeyIndex: 1,
			Output:   "testdata/out.put",
			FillIn:   true,
		},
		false,
	},
	{
		"ios-with-fill-in-and-main-index",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
			FillIn:            true,
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
