package translations

import (
	"testing"
)

var validationCases = []struct {
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
	{
		"all-set-json",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "json",
			KeyIndex: 1,
		},
		true,
	},
	{
		"all-set-json-missing-input",
		Flags{
			Kind:     "json",
			KeyIndex: 1,
		},
		false,
	},
	{
		"all-set-json-missing-key",
		Flags{
			Input: "testdata/input.csv",
			Kind:  "json",
		},
		false,
	},
	{
		"json-with-fill-in-without-main-index",
		Flags{
			Input:    "testdata/input.csv",
			Kind:     "json",
			KeyIndex: 1,
			FillIn:   true,
		},
		false,
	},
	{
		"json-with-fill-in-and-main-index",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "json",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			FillIn:            true,
		},
		true,
	},
	{
		"ios-with-template-valid",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
			FillIn:            true,
			Template:          "testdata/templated-ios-support/file.tmpl",
		},
		true,
	},
	{
		"ios-with-template-invalid",
		Flags{
			Input:             "testdata/input.csv",
			Kind:              "ios",
			KeyIndex:          1,
			DefaultValueIndex: 1,
			Output:            "testdata/out.put",
			FillIn:            true,
			Template:          "does/not/exist.file",
		},
		false,
	},
}

func TestFlagsValidate(t *testing.T) {
	for _, c := range validationCases {
		t.Run(c.name, func(t *testing.T) {
			err := c.flags.validate()
			if c.passes && err != nil {
				t.Errorf("expected to pass, but got %v", err)
			}
			if !c.passes && err == nil {
				t.Errorf("expected to fail, but got %v", err)
			}
		})
	}
}
