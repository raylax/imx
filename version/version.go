package version

import "fmt"

var (
	version = "0.0.1"
	gitSHA  = "Unknown"
	builtOn = "Unknown"
)

func String() string {
	return fmt.Sprintf("Version: %s\nGitSHA: %s\nBuiltOn: %s\n", version, gitSHA, builtOn)
}
