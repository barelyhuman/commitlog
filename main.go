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

	logContainer := logcategory.NewLogsByCategory()
	latestTag, _, err := repo.GetLatestTagFromRepository(r)
	if err != nil {
		return fmt.Errorf("getting tag pairs: %w", err)
	}

	tillLatest := latestTag != nil && latestTag.Hash().String() != ref.Hash().String()

	for _, c := range commits {
		s := c.Hash.String() + " - " + normalizeCommit(c.Message)
		switch k := strings.SplitN(strings.TrimSpace(c.Message), ":", 2)[0]; k {
		case "ci", "fix", "refactor", "feature", "docs":
			logContainer.Add(k, s)
		case "feat":
			logContainer.Add("feature", s)
		default:
			logContainer.Add("", s)
		}

		if nearest, err := isCommitToNearestTag(r, c, tillLatest); err != nil {
			return err
		} else if nearest {
			break
		}
	}

	return logcategory.WriteMarkdown(os.Stdout, logContainer)
}

func isCommitToNearestTag(repository *git.Repository, commit *object.Commit, tillLatest bool) (bool, error) {
	latestTag, previousTag, err := repo.GetLatestTagFromRepository(repository)
	if err != nil {
		return false, fmt.Errorf("get latest tag: %w", err)
	}

	if latestTag == nil {
		return false, nil
	}
	if tillLatest {
		return latestTag.Hash().String() == commit.Hash.String(), nil
	}
	return previousTag.Hash().String() == commit.Hash.String(), nil
}
