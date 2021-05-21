package cmd

import (
	"os"

	commitlogCmd "github.com/barelyhuman/commitlog/cmd/commitlog"
	releaseCmd "github.com/barelyhuman/commitlog/cmd/release"
)

// Init - initializes the commands and also handles the classification
func Init() {
	// Install the needed commands from the cmd package
	releaseCmd.Install()
	commitlogCmd.Install()

	// Classify which command to run, also pass the arguments to the command to parse
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "release":
			{
				releaseCmd.Run(os.Args[2:])
				return
			}
		}
	}

	// Default Command, parse all flags
	commitlogCmd.Run(os.Args)
}
