// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"log"

	clog "github.com/barelyhuman/commitlog/log"
)

func main() {
	var repoPath = flag.String("p", ".", "path to the repository, points to the current working directory by default")
	var startCommit = flag.String("s", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to start collecting commit messages from")
	var endCommit = flag.String("e", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to stop collecting commit message at")
	var inclusionFlags = flag.String("i", "ci,refactor,docs,fix,feat,test,chore,other", "commit types to be includes")
	var skipClassification = flag.Bool("skip", false, "if true will skip trying to classify and just give a list of changes")

	flag.Parse()

	changelog, err := clog.CommitLog(*repoPath, *startCommit, *endCommit, *inclusionFlags, *skipClassification)

	if err.Err != nil {
		log.Fatal(err.Message, err.Err)
	}

	fmt.Println(changelog)

}
