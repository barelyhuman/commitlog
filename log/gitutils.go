// SPDX-License-Identifier: MIT

package commitlog

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	latestTag   *plumbing.Reference
	previousTag *plumbing.Reference
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
	if latestTag == nil || previousTag == nil {
		var err error
		latestTag, previousTag, err = GetLatestTagFromRepository(repo)
		if err != nil {
			log.Fatal("Error getting latest tags from repository")
		}
	}

	ref, err := repo.Head()

	if err != nil {
		log.Fatal("Unable to get repository HEAD:", err)
	}

	tillLatest := latestTag != nil && latestTag.Hash().String() != ref.Hash().String()

	if err != nil {
		log.Fatal("Couldn't get latest tag...", err)
	}

	if latestTag != nil && previousTag != nil {
		if tillLatest {
			return latestTag.Hash() == commit.Hash
		}
		return previousTag.Hash() == commit.Hash

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

// GetCommitFromString - get commit from hash string
func GetCommitFromString(commitString string, repo *git.Repository) *object.Commit {
	if commitString == "" {
		return nil
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(commitString))
	if err != nil {
		log.Fatal("Unable to get Repo head:", err)
	}

	commitRef, err := repo.CommitObject(*hash)
	if err != nil {
		log.Fatal("Unable to get resolve commit:", err)
	}
	return commitRef
}
