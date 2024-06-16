package downloader

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Flags struct {
	Output      string
	Credentials string
	Scopes      string
	DocumentId  string
	Format      string
}

type Client struct {
	outputPath string
	format     string
	documentId string
	scopes     string
	issuer     string
	tokenUri   string
	privateKey any
}

type GoogleCredentials struct {
	PrivateKey string `json:"private_key"`
	Issuer     string `json:"client_email"`
	TokenUri   string `json:"token_uri"`
}

var validFormats = []string{"csv"}

func NewClient(x *Flags) (*Client, error) {
	if err := x.validate(); err != nil {
		return nil, err
	}

	return makeClient(x)
}

func (c *Client) Run() error {
	token, err := c.createAssertionToken()
	if err != nil {
		return err
	}
	accessToken, err := c.fetchAccessToken(token)
	if err != nil {
		return err
	}
	file, err := c.downloadSheet(accessToken)
	if err != nil {
		return err
	}
	if err := c.save(file); err != nil {
		return err
	}
	return nil
}

func (x *Flags) validate() error {
	if len(x.Format) == 0 {
		x.Format = "csv"
	}

	if len(x.Format) == 0 {
		return fmt.Errorf("format not provided")
	}
	if len(x.Output) == 0 {
		return fmt.Errorf("output not provided")
	}
	if len(x.Credentials) == 0 {
		return fmt.Errorf("key not provided")
	}
	if len(x.Scopes) == 0 {
		return fmt.Errorf("scopes not provided")
	}
	if len(x.DocumentId) == 0 {
		return fmt.Errorf("document id not provided")
	}

	for _, v := range validFormats {
		if v != strings.ToLower(x.Format) {
			return fmt.Errorf("invalid format: %s. Valid formats are %v", x.Format, validFormats)
		}
	}

	if _, err := os.Stat(x.Credentials); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Dir(x.Output)); err != nil {
		return err
	}

	return nil
}

func makeClient(x *Flags) (*Client, error) {
	keyBytes, err := os.ReadFile(x.Credentials)
	if err != nil {
		return nil, err
	}

	var gc GoogleCredentials
	err = json.Unmarshal(keyBytes, &gc)
	if err != nil {
		return nil, err
	}

	block, remainer := pem.Decode([]byte(gc.PrivateKey))
	if len(remainer) != 0 {
		return nil, fmt.Errorf("decoding private key had left over bytes")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &Client{
		outputPath: x.Output,
		format:     x.Format,
		documentId: x.DocumentId,
		scopes:     x.Scopes,
		issuer:     gc.Issuer,
		tokenUri:   gc.TokenUri,
		privateKey: key,
	}, nil
}

func (c *Client) createAssertionToken() (string, error) {
	var now = time.Now()
	var token = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   c.issuer,
		"scope": c.scopes,
		"aud":   c.tokenUri,
		"exp":   now.Add(5 * time.Minute).Unix(),
		"iat":   now.Unix(),
	})
	return token.SignedString(c.privateKey)
}

// FIXME: https://golang.cafe/blog/golang-httptest-example.html
func (c *Client) fetchAccessToken(token string) (string, error) {
	// 	response=$(curl --silent --fail -d "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=$assertion" https://oauth2.googleapis.com/token)
	// echo "$response" | jq -r '.access_token'
	return "", nil
}
func (c *Client) downloadSheet(token string) (string, error) {
	// url=$(printf 'https://docs.google.com/spreadsheets/d/%s/export?exportFormat=%s' "$id" "$format")
	// header=$(printf 'Authorization: Bearer %s' "$token")
	return "", nil
}
func (c *Client) save(file string) error {
	// mv?
	return nil
}
