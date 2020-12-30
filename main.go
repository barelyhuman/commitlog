// SPDX-License-Identifier: MIT

package main

import (
	"log"
	"os"
	"strings"

	"github.com/barelyhuman/commitlog/logcategory"
	"github.com/barelyhuman/commitlog/repo"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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

func main() {
	path := os.Args[1]
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal("Error opening Repository: ", err)
	}

	ref, err := r.Head()

	if err != nil {
		log.Fatal("Unable to get repository HEAD:", err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	if err != nil {
		log.Fatal("Unable to get repository commits:", err)
	}

	var commits []*object.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	if err != nil {
		log.Fatal("Error getting commits : ", err)
	}

	logContainer := new(logcategory.LogsByCategory)
	latestTag, _, err := repo.GetLatestTagFromRepository(r)

	if err != nil {
		log.Fatal("Error Getting Tag Pairs", err)
	}

	tillLatest := false

	if latestTag != nil {
		if latestTag.Hash().String() == ref.Hash().String() {
			tillLatest = false
		} else {
			tillLatest = true
		}
	}

	for _, c := range commits {
		switch {
		case strings.Contains(c.Message, "ci:"):
			{
				logContainer.CI = append(logContainer.CI, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "fix:"):
			{
				logContainer.Fix = append(logContainer.Fix, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "refactor:"):
			{
				logContainer.Refactor = append(logContainer.Refactor, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "feat:"):
			{
				logContainer.Feature = append(logContainer.Feature, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "feature:"):
			{
				logContainer.Feature = append(logContainer.Feature, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "docs:"):
			{
				logContainer.Docs = append(logContainer.Docs, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		default:
			{
				logContainer.Other = append(logContainer.Other, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		}

		if isCommitToNearestTag(r, c, tillLatest) {
			break
		}
	}

	logcategory.WriteMarkdown(os.Stdout, logContainer)
}

func isCommitToNearestTag(repository *git.Repository, commit *object.Commit, tillLatest bool) bool {
	latestTag, previousTag, err := repo.GetLatestTagFromRepository(repository)

	if err != nil {
		log.Fatal("Couldn't get latest tag...", err)
	}
	if err != nil {
		log.Fatal("Couldn't access tag...", err)
	}

	if latestTag != nil {
		if tillLatest {
			return latestTag.Hash().String() == commit.Hash.String()
		}
		return previousTag.Hash().String() == commit.Hash.String()

	}
	return false
}
