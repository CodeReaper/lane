package main

import "github.com/codereaper/lane/cmd"
import _ "golang.org/x/crypto/x509roots/fallback"

func main() {
	cmd.Execute()
}
