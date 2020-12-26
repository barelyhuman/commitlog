package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// LogsByCategory - Type to hold logs by each's category
type LogsByCategory struct {
	CI       []string
	FIX      []string
	REFACTOR []string
	FEATURE  []string
	DOCS     []string
	OTHER    []string
}

// GenerateMarkdown - Generate markdown output for the collected commits
func (logContainer *LogsByCategory) GenerateMarkdown() string {
	markDownString := ""

	markDownString += "# Changelog  \n"

	if len(logContainer.CI) > 0 {
		markDownString += "\n\n## CI Changes  \n"

		for _, item := range logContainer.CI {
			markDownString += item + "\n"
		}
	}

	if len(logContainer.FIX) > 0 {
		markDownString += "\n\n## Fixes  \n"
		for _, item := range logContainer.FIX {
			markDownString += item + "\n"
		}
	}

	if len(logContainer.REFACTOR) > 0 {
		markDownString += "\n\n## Performance Fixes  \n"

		for _, item := range logContainer.REFACTOR {
			markDownString += item + "\n"
		}
	}

	if len(logContainer.FEATURE) > 0 {

		markDownString += "\n\n## Feature Fixes  \n"
		for _, item := range logContainer.FEATURE {
			markDownString += item + "\n"
		}
	}

	if len(logContainer.DOCS) > 0 {

		markDownString += "\n\n## Doc Updates  \n"
		for _, item := range logContainer.DOCS {
			markDownString += item + "\n"
		}
	}

	if len(logContainer.OTHER) > 0 {

		markDownString += "\n\n## Other Changes  \n"
		for _, item := range logContainer.OTHER {
			markDownString += item + "\n"
		}
	}

	return markDownString
}

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
	path := os.Args[1]
	r, _ := git.PlainOpen(path)
	ref, _ := r.Head()
	cIter, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	logContainer := new(LogsByCategory)

	_ = cIter.ForEach(func(c *object.Commit) error {
		switch {
		case strings.Contains(c.Message, "ci:"):
			{
				logContainer.CI = append(logContainer.CI, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "fix:"):
			{
				logContainer.FIX = append(logContainer.FIX, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "refactor:"):
			{
				logContainer.REFACTOR = append(logContainer.REFACTOR, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "feature:"):
			{
				logContainer.FEATURE = append(logContainer.FEATURE, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		case strings.Contains(c.Message, "docs:"):
			{
				logContainer.DOCS = append(logContainer.DOCS, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		default:
			{
				logContainer.OTHER = append(logContainer.OTHER, c.Hash.String()+" - "+normalizeCommit(c.Message))
			}
		}
		return nil
	})

	fmt.Println(logContainer.GenerateMarkdown())

}
