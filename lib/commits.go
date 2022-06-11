package lib

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

func GetCommitFromString(repo *git.Repository, commitString string) *object.Commit {
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

func CommitToLog(c *object.Commit) string {
	var commitMsg strings.Builder
	commitMsg.WriteString(c.Hash.String())
	commitMsg.WriteString(" ")
	commitMsg.WriteString(strings.Split(c.Message, "\n")[0])
	return commitMsg.String()
}
