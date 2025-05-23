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
	FillIn            bool
	Template          string
}

func (f *Flags) validate() error {
	if err := f.checkRequiredFields(); err != nil {
		return err
	}

	isIOS := strings.ToLower(f.Kind) == "ios"

	if err := f.checkRequiredFiles(isIOS); err != nil {
		return err
	}

	if isIOS || f.FillIn {
		if f.DefaultValueIndex <= 0 {
			return fmt.Errorf("main index not provided")
		}
	}

	return nil
}

func (f *Flags) checkRequiredFields() error {
	if len(f.Input) == 0 {
		return fmt.Errorf("input not provided")
	}
	if f.KeyIndex <= 0 {
		return fmt.Errorf("index not provided")
	}
	if len(f.Kind) == 0 {
		return fmt.Errorf("kind not provided")
	}

	valid := false
	for _, v := range validKinds {
		if v == strings.ToLower(f.Kind) {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("invalid kind: %s. Valid kinds are %v", f.Kind, validKinds)
	}

	return nil
}

func (f *Flags) checkRequiredFiles(isIOS bool) error {
	if _, err := os.Stat(f.Input); err != nil {
		return err
	}

	if len(f.Template) > 0 {
		if _, err := os.Stat(f.Template); err != nil {
			return err
		}
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
