package version

import "fmt"

var CommitHash string = "unknown"
var Version string = "v0"
var BuildDate string = "unknown"

func String() string {
	return fmt.Sprintf("%s %s %s", Version, CommitHash, BuildDate)
}
