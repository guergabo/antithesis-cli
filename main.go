package main

import (
	"os"

	"github.com/guergabo/antithesis-cli/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		// The error is logged by the CLI library.
		os.Exit(1)
	}
}
