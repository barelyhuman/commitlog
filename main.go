// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
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
	if err := Main(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func Main() error {
	path := os.Args[1]
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("opening repository: %w", err)
	}

	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("get repository HEAD: %w", err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return fmt.Errorf("get repository commits: %w", err)
	}

	var commits []*object.Commit
	if err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	}); err != nil {
		return fmt.Errorf("getting commits: %w", err)
	}

	logContainer := new(logcategory.LogsByCategory)
	latestTag, _, err := repo.GetLatestTagFromRepository(r)
	if err != nil {
		return fmt.Errorf("getting tag pairs: %w", err)
	}

	tillLatest := latestTag != nil && latestTag.Hash().String() != ref.Hash().String()

	for _, c := range commits {
		s := c.Hash.String() + " - " + normalizeCommit(c.Message)
		switch strings.SplitN(strings.TrimSpace(c.Message), ":", 2)[0] {
		case "ci":
			logContainer.CI = append(logContainer.CI, s)
		case "fix":
			logContainer.Fix = append(logContainer.Fix, s)
		case "refactor":
			logContainer.Refactor = append(logContainer.Refactor, s)
		case "feat":
			logContainer.Feature = append(logContainer.Feature, s)
		case "feature":
			logContainer.Feature = append(logContainer.Feature, s)
		case "docs":
			logContainer.Docs = append(logContainer.Docs, s)
		default:
			logContainer.Other = append(logContainer.Other, s)
		}

		if isCommitToNearestTag(r, c, tillLatest) {
			break
		}
	}

	return logcategory.WriteMarkdown(os.Stdout, logContainer)
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
