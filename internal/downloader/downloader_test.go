package downloader

import (
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
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
			Scopes:      "",
			DocumentId:  "",
			Format:      "",
		},
		false,
	},
	{
		"all-set",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		true,
	},
	{
		"all-set-multiple-scopes",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid email custom",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		true,
	},
	{
		"all-set-but-output",
		Flags{
			Output:      "",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
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
			Scopes:      "openid",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-scopes",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-documentid",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
			DocumentId:  "",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-but-format",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
			DocumentId:  "1234567890",
			Format:      "",
		},
		true,
	},
	{
		"all-set-incorrect-format",
		Flags{
			Output:      "testdata/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
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
			Scopes:      "openid",
			DocumentId:  "1234567890",
			Format:      "csv",
		},
		false,
	},
	{
		"all-set-missing-output-directory",
		Flags{
			Output:      "not-here/out.csv",
			Credentials: "testdata/googleapi.json",
			Scopes:      "openid",
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

func TestAssertionToken(t *testing.T) {
	config := Flags{
		Output:      "testdata/out.csv",
		Credentials: "testdata/googleapi.json",
		Scopes:      "openid",
		DocumentId:  "1234567890",
		Format:      "csv",
	}
	client, err := NewClient(&config)
	if err != nil {
		t.Error(err)
	}
	createdToken, err := client.createAssertionToken()
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, createdToken)

	privateKey, ok := client.privateKey.(*rsa.PrivateKey)
	assert.True(t, ok)

	parsedToken, err := jwt.Parse(createdToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)
}
