// Command ledgerly unifies bank CSV exports into one clean, categorized format.
package main

import (
	"os"

	"github.com/LukaZagar/ledgerly/internal/cli"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		os.Exit(1)
	}
}
