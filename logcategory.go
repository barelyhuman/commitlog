// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"
)

// logContainer - Container of log strings
type logContainer struct {
	include bool
	commits []string
}

// logsByCategory - Type to hold logs by each's category
// to be left as ALLCAPS to be considered as symbols instead of selectors
type logsByCategory struct {
	CI           logContainer
	FIX          logContainer
	REFACTOR     logContainer
	FEATURE      logContainer
	DOCS         logContainer
	CHORE        logContainer
	TEST         logContainer
	OTHER        logContainer
	UNCLASSIFIED logContainer
}

// Setup - Initialize all Log Containers
func (logs logsByCategory) Setup() {
	logs.CI.include = true
	logs.FIX.include = true
	logs.REFACTOR.include = true
	logs.FEATURE.include = true
	logs.DOCS.include = true
	logs.CHORE.include = true
	logs.TEST.include = true
	logs.OTHER.include = true
	logs.UNCLASSIFIED.include = true
}

// printLog - loops through the collected logs to write them to string builder
func (logs logContainer) printLog(out *strings.Builder, title string, skipped bool) {
	if !logs.include {
		return
	}
	if len(logs.commits) > 0 {
		if !skipped {
			out.WriteString(fmt.Sprintf("\n\n## %s  \n", title))
		}
		for _, item := range logs.commits {
			out.WriteString(item + "\n")
		}
	}
}

// ToMarkdown - Generate markdown output for the collected commits
func (logs *logsByCategory) ToMarkdown(skipped bool) string {
	var markdownString strings.Builder

	markdownString.WriteString("# Changelog \n")

	if !skipped {
		logs.FEATURE.printLog(&markdownString, "Features", skipped)
		logs.FIX.printLog(&markdownString, "Fixes", skipped)
		logs.REFACTOR.printLog(&markdownString, "Performance", skipped)
		logs.CI.printLog(&markdownString, "CI", skipped)
		logs.DOCS.printLog(&markdownString, "Docs", skipped)
		logs.CHORE.printLog(&markdownString, "Chores", skipped)
		logs.TEST.printLog(&markdownString, "Tests", skipped)
		logs.OTHER.printLog(&markdownString, "Other Changes", skipped)
	} else {
		logs.UNCLASSIFIED.include = true
		logs.UNCLASSIFIED.printLog(&markdownString, "Unclassified Changes", skipped)
	}

	return markdownString.String()
}

// AddCommit - Add a commit to the needed logContainer based on skip and include flag
func (logs *logsByCategory) AddCommit(key, commitHash string, skip bool) {
	var addCommitToContainer *logContainer
	switch key {
	case "ci":
		if logs.CI.include && !skip {
			addCommitToContainer = &logs.CI
		} else if skip && logs.CI.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "fix":
		if logs.FIX.include && !skip {
			addCommitToContainer = &logs.FIX
		} else if skip && logs.FIX.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "refactor":
		if logs.REFACTOR.include && !skip {
			addCommitToContainer = &logs.REFACTOR
		} else if skip && logs.REFACTOR.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "feat", "feature":
		if logs.FEATURE.include && !skip {
			addCommitToContainer = &logs.FEATURE
		} else if skip && logs.FEATURE.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "docs":
		if logs.DOCS.include && !skip {
			addCommitToContainer = &logs.DOCS
		} else if skip && logs.DOCS.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "test":
		if logs.TEST.include && !skip {
			addCommitToContainer = &logs.TEST
		} else if skip && logs.TEST.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	case "chore":
		if logs.CHORE.include && !skip {
			addCommitToContainer = &logs.CHORE
		} else if skip && logs.CHORE.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	default:
		if logs.OTHER.include && !skip {
			addCommitToContainer = &logs.OTHER
		} else if skip && logs.OTHER.include {
			addCommitToContainer = &logs.UNCLASSIFIED
		}
	}

	if addCommitToContainer != nil {
		addCommitToContainer.commits = append(addCommitToContainer.commits, commitHash)
	}
}
