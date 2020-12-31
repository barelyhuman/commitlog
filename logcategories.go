package main

import (
	"fmt"
	"strings"
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

// printCategory will output all items inside a Log slice and a title
func printCategory(output *strings.Builder, title string, logs []string) {
	if len(logs) > 0 {
		output.WriteString(fmt.Sprintf("\n\n## %s  \n", title))
		for _, item := range logs {
			output.WriteString(item + "\n")
		}
	}
}

// GenerateMarkdown - Generate markdown output for the collected commits
func (logContainer *LogsByCategory) GenerateMarkdown() string {
	var output strings.Builder
	output.WriteString("# Changelog \n")

	printCategory(&output, "CI Changes", logContainer.CI)
	printCategory(&output, "Fixes", logContainer.FIX)
	printCategory(&output, "Performance Fixes", logContainer.REFACTOR)
	printCategory(&output, "Feature Fixes", logContainer.FEATURE)
	printCategory(&output, "Doc Updates", logContainer.DOCS)
	printCategory(&output, "Other Changes", logContainer.OTHER)

	return output.String()
}

// AddCommitLog will take a commitHash and a commitMessage and append them to their
// apropriate Slice
func (logContainer *LogsByCategory) AddCommitLog(commitHash, commitMessage string) {
	message := fmt.Sprintf("%s - %s", commitHash, normalizeCommit(commitMessage))

	switch {
	case strings.Contains(commitMessage, "ci:"):
		logContainer.CI = append(logContainer.CI, message)
	case strings.Contains(commitMessage, "fix:"):
		logContainer.FIX = append(logContainer.FIX, message)
	case strings.Contains(commitMessage, "refactor:"):
		logContainer.REFACTOR = append(logContainer.REFACTOR, message)
	// Golang Switch allows multiple values in cases with , separation
	case strings.Contains(commitMessage, "feat:"), strings.Contains(commitMessage, "feature:"):
		logContainer.FEATURE = append(logContainer.FEATURE, message)
	case strings.Contains(commitMessage, "docs:"):
		logContainer.DOCS = append(logContainer.DOCS, message)

	default:
		logContainer.OTHER = append(logContainer.OTHER, message)
	}
}
