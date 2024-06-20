package translations

import (
	"context"
	"encoding/csv"
	"log"
	"os"
)

var androidKind = "android"
var iosKind = "ios"
var jsonKind = "json"

var validKinds = []string{
	iosKind,
	androidKind,
	jsonKind,
}

type Generator struct {
	flags *Flags
}

func NewGenerator(x *Flags) *Generator {
	return &Generator{
		flags: x,
	}
}

func (g *Generator) Generate(ctx context.Context) error {
	if err := g.flags.validate(); err != nil {
		return err
	}

	r, err := os.Open(g.flags.Input)
	if err != nil {
		return err
	}

	reader := csv.NewReader(r)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) > 0 {
		records = records[1:] // skips header row
	}

	switch g.flags.Kind {
	case androidKind:
		return writeAndroid(records, g.flags.Configurations)
	case iosKind:
		return writeIos(records, g.flags.Configurations)
	case jsonKind:
		return writeJson(records, g.flags.Configurations)
	default:
		log.Fatalf("found unknown kind: %v", g.flags.Kind)
	}

	return nil
}

func writeAndroid(records [][]string, configurations []string) error {
	panic("unimplemented")
}

func writeIos(records [][]string, configurations []string) error {
	panic("unimplemented")
}

func writeJson(records [][]string, configurations []string) error {
	panic("unimplemented")
}
