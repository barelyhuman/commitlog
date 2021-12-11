// SPDX-License-Identifier: MIT

package commitlog

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func IsHashATag(currentRepository *git.Repository, hash plumbing.Hash) bool {
	isTag := false
	tags, _ := currentRepository.Tags()
	tags.ForEach(func(tagRef *plumbing.Reference) error {
		revHash, err := currentRepository.ResolveRevision(plumbing.Revision(tagRef.Name()))
		if err != nil {
			return err
		}
		if *revHash == hash {
			isTag = true
		}
		return nil
	})
	return isTag
}

// normalizeCommit - reduces the commit message to the first line and ignore the description text of the commit
func normalizeCommit(commitMessage string, key string) string {
	var message string
	for i, msg := range strings.Split(commitMessage, "\n") {
		if i == 0 {
			message = msg
			break
		}
	}

	removedPrefix := strings.TrimPrefix(strings.TrimSuffix(message, "\n"), key)
	return strings.TrimSpace(strings.TrimSuffix(removedPrefix, "\n"))
}

// OpenRepository - open the git repository and return repository reference
func OpenRepository(path string) *git.Repository {
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
