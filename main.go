// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
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

type commitTypeInclusions [][]byte

func main() {
	var repoPath = flag.String("p", ".", "path to the repository, points to the current working directory by default")
	var startCommit = flag.String("s", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to start collecting commit messages from")
	var endCommit = flag.String("e", "", "commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) \n to stop collecting commit message at")
	var inclusionFlags = flag.String("i", "ci,refactor,docs,fix,feat,test,chore,other", "commit types to be includes")
	var skipClassification = flag.Bool("skip", false, "if true will skip trying to classify and just give a list of changes")

	flag.Parse()

	err := CommitLog(*repoPath, *startCommit, *endCommit, *inclusionFlags, *skipClassification)

	if err.err != nil {
		log.Fatal(err.message, err.err)
	}
}

// CommitLog - Generate commit log
func CommitLog(path string, startCommitString string, endCommitString string, inclusionFlags string, skipClassification bool) ErrMessage {
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

	var inclusions commitTypeInclusions = bytes.SplitN([]byte(inclusionFlags), []byte(","), -1)

	logContainer := logsByCategory{}

	logContainer.Setup()

	logContainer.CI.include = inclusions.checkInclusion("ci")
	logContainer.FIX.include = inclusions.checkInclusion("fix")
	logContainer.REFACTOR.include = inclusions.checkInclusion("refactor")
	logContainer.FEATURE.include = inclusions.checkInclusion("feat")
	logContainer.DOCS.include = inclusions.checkInclusion("docs")
	logContainer.CHORE.include = inclusions.checkInclusion("chore")
	logContainer.TEST.include = inclusions.checkInclusion("test")
	logContainer.OTHER.include = inclusions.checkInclusion("other")

	for _, c := range commits {
		normalizedHash := c.Hash.String() + " - " + normalizeCommit(c.Message)
		key := strings.SplitN(strings.TrimSpace(c.Message), ":", 2)[0]

		logContainer.AddCommit(key, normalizedHash, skipClassification)

		if endHash == nil && isCommitToNearestTag(currentRepository, c) {
			break
		} else if endHash != nil && c.Hash == endHash.Hash {
			break
		}
	}
	fmt.Println(logContainer.ToMarkdown(skipClassification))

	return ErrMessage{}
}

func (inclusions *commitTypeInclusions) checkInclusion(flagToCheck string) bool {
	if inclusions != nil {
		for _, flag := range *inclusions {
			if string(flag) == flagToCheck {
				return true
			}
		}
	}
	return false
}
