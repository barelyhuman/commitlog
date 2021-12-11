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
	clogCmd = flag.NewFlagSet("commitlog", flag.ExitOnError)
	repoPath = clogCmd.String("p", ".", "path to the repository, points to the current working directory by default")
	startCommit = clogCmd.String("s", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to start collecting commit messages from")
	endCommit = clogCmd.String("e", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to stop collecting commit message at")
	inclusionFlags = clogCmd.String("i", clog.SupportedKeys, "commit types to be includes")
	skipClassification = clogCmd.Bool("skip", false, "if true will skip trying to classify and just give a list of changes")
}

// Run - execute the command
func Run(args []string) {

	err := clogCmd.Parse(args)

	if err != nil {
		log.Fatalln(err)
	}

	currentRepository := clog.OpenRepository(*repoPath)

	changelog, clogErr := clog.CommitLog(currentRepository, *startCommit, *endCommit, *inclusionFlags, *skipClassification)

	if clogErr.Err != nil {
		log.Fatal(clogErr.Message, clogErr.Err)
	}

	fmt.Println(changelog)
}
