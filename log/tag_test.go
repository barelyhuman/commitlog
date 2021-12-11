package commitlog

import (
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getTagOptions(message string) *git.CreateTagOptions {
	return &git.CreateTagOptions{
		Message: message,
		Tagger: &object.Signature{
			Name:  "Test",
			Email: "test@reaper.im",
			When:  time.Now(),
		},
	}

}

func TestCommitLogSingleTag(t *testing.T) {
	secondCommit := expectedCommits[1]
	acceptedCommits := expectedCommits[2:]

	t.Log("Commits:", expectedCommits)
	t.Log("Tagged:", secondCommit)

	hash, err := repo.ResolveRevision(plumbing.Revision(secondCommit))
	bail(err)

	_, err = repo.CreateTag("0.0.0", *hash, getTagOptions("0.0.0"))
	bail(err)

	log, _ := CommitLog(repo, "", "", SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	if strings.Contains(log, expectedCommits[0]) {
		t.Fail()
	}

	for _, commit := range acceptedCommits {
		if !strings.Contains(log, commit) {
			t.Fail()
		}
	}

	t.Log(log)

	// clean-up
	bail(repo.DeleteTag("0.0.0"))

}

// Test with 2 tags, one on the second commit and one on the 2nd last commit,
// should only have the last commit in the log
func TestCommitLogDualTag(t *testing.T) {
	secondCommit := expectedCommits[1]
	secondLastCommit := expectedCommits[len(expectedCommits)-2]
	acceptedCommit := expectedCommits[len(expectedCommits)-1]

	t.Log("Commits:", expectedCommits)
	t.Log("Tagged:", secondCommit, secondLastCommit)

	secondHash, err := repo.ResolveRevision(plumbing.Revision(secondCommit))
	bail(err)

	secondLastHash, err := repo.ResolveRevision(plumbing.Revision(secondLastCommit))
	bail(err)

	_, err = repo.CreateTag("0.0.0", *secondHash, getTagOptions("0.0.0"))
	bail(err)

	_, err = repo.CreateTag("0.0.1", *secondLastHash, getTagOptions("0.0.1"))
	bail(err)

	log, _ := CommitLog(repo, "", "", SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	for _, commit := range expectedCommits {
		if commit == acceptedCommit {
			if !strings.Contains(log, acceptedCommit) {
				t.Fail()
			}
		} else {
			if strings.Contains(log, commit) {
				t.Fail()
			}
		}

	}

	t.Log(log)

	// clean-up
	bail(repo.DeleteTag("0.0.0"))
	bail(repo.DeleteTag("0.0.1"))
}

// Test with 2 tags, one on the second commit and one on the last commit,
// should give all commits till the 1st tag
func TestCommitLogHeadTag(t *testing.T) {
	secondCommit := expectedCommits[1]
	lastCommit := expectedCommits[len(expectedCommits)-1]

	t.Log("Commits:", expectedCommits)
	t.Log("Tagged:", secondCommit, lastCommit)

	secondHash, err := repo.ResolveRevision(plumbing.Revision(secondCommit))
	bail(err)

	lastHash, err := repo.ResolveRevision(plumbing.Revision(lastCommit))
	bail(err)

	_, err = repo.CreateTag("0.0.0", *secondHash, getTagOptions("0.0.0"))
	bail(err)

	_, err = repo.CreateTag("0.0.1", *lastHash, getTagOptions("0.0.1"))
	bail(err)

	log, _ := CommitLog(repo, "", "", SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	for index, commit := range expectedCommits {
		if index <= 1 {
			if strings.Contains(log, commit) {
				t.Fail()
			}
		}
		if index > 1 && !strings.Contains(log, commit) {
			t.Fail()
		}

	}

	t.Log(log)

	// clean-up
	bail(repo.DeleteTag("0.0.0"))
	bail(repo.DeleteTag("0.0.1"))
}
