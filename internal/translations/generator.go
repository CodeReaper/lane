package translations

import (
	"context"
)

var androidKind = "android"
var iosKind = "ios"
var jsonKind = "json"

var validKinds = []string{
	iosKind,
	androidKind,
	jsonKind,
}

func Generate(ctx context.Context, flags *Flags, configurations []string) error {
	err := flags.validate()
	if err != nil {
		return err
	}

	files, err := newLanguageFiles(flags, configurations)
	if err != nil {
		return err
	}

	translations, err := newTranslations(flags.Input)
	if err != nil {
		return err
	}

	for _, f := range files {
		err := f.Write(translations)
		if err != nil {
			return err
		}
	}

	return nil
}
