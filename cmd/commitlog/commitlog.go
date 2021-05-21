package commitlog

import (
	"flag"
	"fmt"
	"log"

	clog "github.com/barelyhuman/commitlog/log"
)

var clogCmd *flag.FlagSet
var repoPath *string
var startCommit *string
var endCommit *string
var inclusionFlags *string
var skipClassification *bool

// Install - add flags and other options
func Install() {
	clogCmd = flag.NewFlagSet("release", flag.ExitOnError)
	repoPath = clogCmd.String("p", ".", "path to the repository, points to the current working directory by default")
	startCommit = clogCmd.String("s", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to start collecting commit messages from")
	endCommit = clogCmd.String("e", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to stop collecting commit message at")
	inclusionFlags = clogCmd.String("i", clog.SupportedKeys, "commit types to be includes")
	skipClassification = clogCmd.Bool("skip", false, "if true will skip trying to classify and just give a list of changes")
}

// Run - execute the command
func Run(args []string) {
	clogCmd.Parse(args)
	changelog, err := clog.CommitLog(*repoPath, *startCommit, *endCommit, *inclusionFlags, *skipClassification)

	if err.Err != nil {
		log.Fatal(err.Message, err.Err)
	}

	fmt.Println(changelog)
}
