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
func (container logContainer) printLog(out *strings.Builder, title string, skipped bool) {
	if !container.include {
		return
	}
	if len(container.commits) > 0 {
		if !skipped {
			out.WriteString(fmt.Sprintf("\n\n## %s  \n", title))
		}
		for _, item := range container.commits {
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
	addCommitToContainer := logs.findContainerByKey(key)
	if !addCommitToContainer.canAddToContainer(skip) {
		addCommitToContainer = &logs.UNCLASSIFIED
	}
	if addCommitToContainer != nil {
		addCommitToContainer.commits = append(addCommitToContainer.commits, commitHash)
	}
}

func (logs *logsByCategory) findContainerByKey(key string) *logContainer {
	switch key {
	case "ci":
		return &logs.CI
	case "fix":
		return &logs.FIX
	case "refactor":
		return &logs.REFACTOR
	case "feat", "feature":
		return &logs.FEATURE
	case "docs":
		return &logs.DOCS
	case "test":
		return &logs.TEST
	case "chore":
		return &logs.CHORE
	default:
		return &logs.OTHER
	}
}

func (container *logContainer) canAddToContainer(skip bool) bool {
	if container.include && !skip {
		return true
	} else if skip && container.include {
		return false
	}
	return true
}
