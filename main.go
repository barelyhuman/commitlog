// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ErrMessage - simple interface around error with a custom message
type ErrMessage struct {
	message string
	err     error
}

func main() {
	var repoPath = flag.String("p", ".", "path to the repository, points to the current working directory by default")
	var startCommit = flag.String("s", "", "commit hash string start collecting commit messages from")
	var endCommit = flag.String("e", "", "commit hash string to stop collecting commit message at")

	flag.Parse()

	path := repoPath

	err := CommitLog(*path, *startCommit, *endCommit)

	if err.err != nil {
		log.Fatal(err.message, err.err)
	}
}

// CommitLog - Generate commit log
func CommitLog(path string, startCommitString string, endCommitString string) ErrMessage {
	currentRepository := openRepository(path)

	baseCommitReference, err := currentRepository.Head()
	var startHash, endHash *object.Commit
	var cIter object.CommitIter

	if err != nil {
		return ErrMessage{"Unable to get repository HEAD:", err}
	}

	if startCommitString != "" {
		startHash = GetCommitFromString(startCommitString, currentRepository)
	}

	if endCommitString != "" {
		endHash = GetCommitFromString(endCommitString, currentRepository)
	}

	if startHash != nil {
		cIter, err = currentRepository.Log(&git.LogOptions{From: startHash.Hash})
	} else {
		cIter, err = currentRepository.Log(&git.LogOptions{From: baseCommitReference.Hash()})
	}

	if err != nil {
		return ErrMessage{"Unable to get repository commits:", err}
	}

	var commits []*object.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	if err != nil {
		return ErrMessage{"Error getting commits : ", err}
	}

	logContainer := logsByCategory{}

	for _, c := range commits {
		normalizedHash := c.Hash.String() + " - " + normalizeCommit(c.Message)
		switch strings.SplitN(strings.TrimSpace(c.Message), ":", 2)[0] {
		case "ci":
			logContainer.CI = append(logContainer.CI, normalizedHash)
		case "fix":
			logContainer.FIX = append(logContainer.FIX, normalizedHash)
		case "refactor":
			logContainer.REFACTOR = append(logContainer.REFACTOR, normalizedHash)
		case "feat", "feature":
			logContainer.FEATURE = append(logContainer.FEATURE, normalizedHash)
		case "docs":
			logContainer.DOCS = append(logContainer.DOCS, normalizedHash)
		case "test":
			logContainer.TEST = append(logContainer.TEST, normalizedHash)
		case "chore":
			logContainer.CHORE = append(logContainer.CHORE, normalizedHash)
		default:
			logContainer.OTHER = append(logContainer.OTHER, normalizedHash)
		}

		if endHash == nil && isCommitToNearestTag(currentRepository, c) {
			break
		} else if endHash != nil && c.Hash == endHash.Hash {
			break
		}
	}
	fmt.Println(logContainer.ToMarkdown())

	return ErrMessage{}
}
