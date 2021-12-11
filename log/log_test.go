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

func setup() *git.Repository {
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

var repo *git.Repository = setup()

func TestCommitLogDefault(t *testing.T) {

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

func TestCommitLogStartHash(t *testing.T) {

	expectedCommitsLen := len(expectedCommits)
	startCommitHash := expectedCommits[expectedCommitsLen-2]
	lastCommit := expectedCommits[expectedCommitsLen-1]
	acceptedCommitHashes := expectedCommits[0 : expectedCommitsLen-1]

	t.Log("Commits: ", expectedCommits)
	t.Log("Start At:", startCommitHash)

	log, _ := CommitLog(repo, startCommitHash, "", SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	// should have all commits except the last one
	if strings.Contains(log, lastCommit) {
		t.Fail()
	}

	for _, commitHash := range acceptedCommitHashes {
		if !strings.Contains(log, commitHash) {
			t.Log("Failed at:", commitHash)
			t.Fail()
		}
	}

	t.Log("\n", log)
}

func TestCommitLogEndHash(t *testing.T) {

	endCommitHash := expectedCommits[1]
	firstCommit := expectedCommits[0]
	acceptedCommitHashes := expectedCommits[2:]

	t.Log("Commits: ", expectedCommits)
	t.Log("End At:", endCommitHash)

	log, _ := CommitLog(repo, "", endCommitHash, SupportedKeys, true)
	if log == "" {
		t.Fail()
	}

	// should have all commits except the first one
	if strings.Contains(log, firstCommit) {
		t.Fail()
	}

	for _, commitHash := range acceptedCommitHashes {
		if !strings.Contains(log, commitHash) {
			t.Log("Failed at:", commitHash)
			t.Fail()
		}
	}

	t.Log("\n", log)
}
