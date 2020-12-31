// SPDX-License-Identifier: MIT

package main

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetLatestTagFromRepository - Get the latest Tag reference from the repo
func GetLatestTagFromRepository(repository *git.Repository) (*plumbing.Reference, *plumbing.Reference, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return nil, nil, err
	}

	var latestTagCommit *object.Commit
	var latestTagName *plumbing.Reference
	var previousTag *plumbing.Reference
	var previousTagReturn *plumbing.Reference

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())

		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef
			previousTagReturn = previousTag
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef
			previousTagReturn = previousTag
		}

		previousTag = tagRef

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return latestTagName, previousTagReturn, nil
}

// isCommitToNearestTag -  go through git revisions to find the latest tag and the nearest next tag
func isCommitToNearestTag(repo *git.Repository, commit *object.Commit) bool {
	latestTag, previousTag, err := GetLatestTagFromRepository(repo)

	ref, err := repo.Head()

	if err != nil {
		log.Fatal("Unable to get repository HEAD:", err)
	}

	tillLatest := latestTag != nil && latestTag.Hash().String() != ref.Hash().String()

	if err != nil {
		log.Fatal("Couldn't get latest tag...", err)
	}

	if latestTag != nil {
		if tillLatest {
			return latestTag.Hash().String() == commit.Hash.String()
		}
		return previousTag.Hash().String() == commit.Hash.String()

	}
	return false
}

// normalizeCommit - reduces the commit message to the first line and ignore the description text of the commit
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

// openRepository - open the git repository and return repository reference
func openRepository(path string) *git.Repository {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal("Error opening Repository: ", err)
	}

	return r
}
