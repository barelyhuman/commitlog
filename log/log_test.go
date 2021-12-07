package commitlog

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var testCommits []string = []string{
	"fix: fix commit",
	"feat: feat commit",
	"docs: doc update commit",
	"chore: chore commit",
	"other commit",
}

var expectedCommits []string

func bail(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func setup(t *testing.T) *git.Repository {
	var fs = memfs.New()
	repo, _ := git.Init(memory.NewStorage(), fs)
	wt, err := repo.Worktree()
	bail(err)

	for _, testCommitMsg := range testCommits {
		commit, err := wt.Commit(testCommitMsg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Reaper",
				Email: "ahoy@barelyhuman.dev",
				When:  time.Now(),
			},
		})
		bail(err)
		expectedCommits = append(expectedCommits, commit.String())
	}

	return repo
}

func TestCommitLogDefault(t *testing.T) {
	repo := setup(t)

	log, _ := CommitLog(repo, "", "", SupportedKeys, false)
	if log == "" {
		t.Fail()
	}

	for _, commit := range expectedCommits {
		if !strings.Contains(log, commit) {
			t.Fail()
		}
	}

	t.Log(log)

}

func TestCommitLogSkipped(t *testing.T) {
	repo := setup(t)

	log, _ := CommitLog(repo, "", "", SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	for _, commit := range expectedCommits {
		// Shouldn't contain classification headings
		if strings.Contains(log, "##") {
			t.Fail()
		}

		if !strings.Contains(log, commit) {
			t.Fail()
		}
	}

	t.Log(log)
}

func TestCommitLogInclusions(t *testing.T) {
	repo := setup(t)

	// include only feature commits
	log, _ := CommitLog(repo, "", "", "feat", true)
	if log == "" {
		t.Fail()
	}

	ignoredHeadings := []string{
		"## Fixes",
		"## Performance",
		"## CI",
		"## Docs",
		"## Chores",
		"## Tests",
		"## Other Changes",
	}

	for _, heading := range ignoredHeadings {
		if strings.Contains(log, heading) {
			t.Fail()
		}
	}

	t.Log(log)
}

// TODO:
// - Tests for checking between tags
// - Variation of the above to check between 2 tags
// - Another variation where one tag points to the head of the repo
// - Test for checking start and end commit hashes passed as parameters
