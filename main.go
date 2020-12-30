package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	var path string
	// Read user input
	flag.StringVar(&path, "path", "", "A filepath to a folder containing a github repository")
	// Parse Flags
	flag.Parse()

	// Make sure user has inserted the needed flags
	if path == "" {
		flag.Usage()
		os.Exit(0)
	}

	repo, err := Open(path)
	if err != nil {
		log.Fatal(err)
	}

	commits, err := repo.GetCommits()
	if err != nil {
		log.Fatal(err)
	}

	logContainer := new(LogsByCategory)

	// we no longer need to fetch latestTag here to compare tillLatest.

	// itterate all commits and add them to the log based on hash and Message
	for _, c := range commits {

		logContainer.AddCommitLog(c.Hash.String(), c.Message)

		nearestTag, err := repo.IsCommitNearest(c)
		if err != nil {
			log.Fatal(err)
		}
		if nearestTag {
			break
		}
	}

	fmt.Println(logContainer.GenerateMarkdown())

}

func normalizeCommit(commitMessage string) string {
	var message string
	for i, msg := range strings.Split(commitMessage, "\n") {
		if i == 0 {
			message = msg
			break
		}
	}
	return strings.TrimSuffix(message, "\n")
}
