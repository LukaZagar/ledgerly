// Command ledgerly unifies bank CSV exports into one clean, categorized format.
package main

import "fmt"

// version is overridden at build time via -ldflags.
var version = "dev"

func main() {
	fmt.Printf("ledgerly %s\n", version)
}
