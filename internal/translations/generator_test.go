package translations

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	expected := `<resources>
    <string name="something">Something</string>
    <string name="something_with_arguments">Something with %1$s and %2$s</string>
</resources>`

	flags := Flags{
		Input:          "testdata/input.csv",
		Kind:           "android",
		Index:          1,
		Configurations: []string{"3 build/en.xml"},
	}

	generator := NewGenerator(&flags)
	err := generator.Generate(context.Background())

	assert.Nil(t, err)
	b, err := os.ReadFile("build/en.xml")
	assert.Nil(t, err)
	assert.Equal(t, expected, string(b))
}
