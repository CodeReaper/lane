package translations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Flags struct {
	Input             string
	Kind              string
	KeyIndex          int
	DefaultValueIndex int
	Output            string
}

func (f *Flags) validate() error {
	if len(f.Input) == 0 {
		return fmt.Errorf("input not provided")
	}
	if len(f.Kind) == 0 {
		return fmt.Errorf("kind not provided")
	}

	isIOS := false

	validKind := false
	for _, v := range validKinds {
		if !validKind && v == strings.ToLower(f.Kind) {
			validKind = true
		}
		isIOS = isIOS || strings.ToLower(f.Kind) == "ios"
	}
	if !validKind {
		return fmt.Errorf("invalid kind: %s. Valid kinds are %v", f.Kind, validKinds)
	}

	if _, err := os.Stat(f.Input); err != nil {
		return err
	}

	if isIOS {
		if len(f.Output) == 0 {
			return fmt.Errorf("output not provided")
		}
		if _, err := os.Stat(filepath.Dir(f.Output)); err != nil {
			return err
		}
	}

	return nil
}
