package release

import (
	"flag"
	"fmt"
)

var releaseCmd *flag.FlagSet

// Install - add flags and other options
func Install() {
	releaseCmd = flag.NewFlagSet("release", flag.ExitOnError)
	releaseCmd.Bool("-major", false, "If release is a major one, will increment the x.0.0 ")
	releaseCmd.Bool("-minor", false, "If release is a minor one, will increment the 0.x.0 ")
	releaseCmd.Bool("-patch", false, "If release is a patch, will increment the 0.0.x ")
	releaseCmd.Bool("-beta", false, "If the release is a beta, to add/increment tag by `-beta.x`")
	releaseCmd.String("-tag", "", "The Tag to be taken as base")
}

// Run - execute the command
func Run(args []string) {
	releaseCmd.Parse(args)
	fmt.Println("Note: The release command is not yet implemented")
}
