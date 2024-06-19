package downloader

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
			Output:      "",
			Credentials: "",
			DocumentId:  "",
			Format:      "",
		},
		false,
	},
	{
		"all-set",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		true,
	},
	{
		"all-set-multiple-scopes",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		true,
	},
	{
		"all-set-but-output",
		Flags{
			Output:      "",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-key",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-documentid",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-format",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "",
		},
		true,
	},
	{
		"all-set-incorrect-format",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "unknown",
		},
		false,
	},
	{
		"all-set-missing-key",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/not-here.json",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-missing-output-directory",
		Flags{
			Output:      "not-here/out.csv",
			Credentials: "testdata/empty.json",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
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
