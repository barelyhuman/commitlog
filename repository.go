package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Repository is a struct that holds an Opened github repository and the reference
type Repository struct {
	path string
	repo *git.Repository
	ref  *plumbing.Reference
}

// Open will open up a git repository
// and load the Head reference
func Open(path string) (*Repository, error) {
	// Open the github repository
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", "Error opening Repository: ", err)
	}

	// Grab the Git HEAD reference
	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("%v:%w", "Unable to get repository HEAD:", err)
	}
	return &Repository{
		path: path,
		repo: r,
		ref:  ref,
	}, nil
}

// GetCommits will extract commit history from the git repository
func (r *Repository) GetCommits() ([]*object.Commit, error) {
	// Get the Commit history
	cIter, err := r.repo.Log(&git.LogOptions{From: r.ref.Hash()})

	if err != nil {
		return nil, fmt.Errorf("%v:%w", "Unable to get repository commits:", err)
	}

	var commits []*object.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%v:%w", "Error getting commits :", err)
	}
	return commits, nil
}

// GetLatestTag is used to check the latestTag
// it will return a reference to the LatestTag and the PreviousTag
func (r *Repository) GetLatestTag() (*plumbing.Reference, *plumbing.Reference, error) {
	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, nil, err
	}

	var latestTagCommit *object.Commit
	var latestTagName *plumbing.Reference
	var previousTag *plumbing.Reference
	var previousTagReturn *plumbing.Reference

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())

		tagCommitHash, err := r.repo.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := r.repo.CommitObject(*tagCommitHash)
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

// IsCommitNearest will check if a commit tag Hash is equal to the current repository HEAD tag
// If the Hashes matches, it will return true
func (r *Repository) IsCommitNearest(commit *object.Commit) (bool, error) {
	latestTag, previousTag, err := r.GetLatestTag()

	if err != nil {
		return false, fmt.Errorf("%v:%w", "Couldn't get latest tag...", err)
	}

	if latestTag != nil {
		// Compare latest tag Hash with the repository HEAD hash
		// Hash() returns a Slice which can be compared without converting to string
		// Noticed by /OfficialTomCruise on reddit comments
		if latestTag.Hash() == r.ref.Hash() {
			return true, nil
		}
		return previousTag.Hash() == commit.Hash, nil

	}
	return false, nil
}
