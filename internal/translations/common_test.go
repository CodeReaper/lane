package translations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func equalFiles(t *testing.T, expectedPath string, actualPath string) bool {
	expected, err := os.ReadFile(expectedPath)
	if !assert.NoError(t, err) {
		return false
	}

	actual, err := os.ReadFile(actualPath)
	if !assert.NoError(t, err) {
		return false
	}

	return assert.Equal(t, string(expected), string(actual))
}
