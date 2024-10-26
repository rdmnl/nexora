package version

import "fmt"

var (
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "none"
)

func PrintVersion() {
	fmt.Printf("Nexora version: %s, Build date: %s, Commit: %s\n", Version, BuildDate, Commit)
}
