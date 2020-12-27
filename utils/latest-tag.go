package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetLatestTagFromRepository - Get the latest Tag reference from the repo
func GetLatestTagFromRepository(repository *git.Repository) (*plumbing.Reference, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return nil, err
	}

	var latestTagCommit *object.Commit
	var latestTagName *plumbing.Reference
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
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return latestTagName, nil
}
