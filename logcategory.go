// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"
)

// logContainer - Container of log strings
type logContainer []string

// logsByCategory - Type to hold logs by each's category
// to be left as ALLCAPS to be considered as symbols instead of selectors
type logsByCategory struct {
	CI       logContainer
	FIX      logContainer
	REFACTOR logContainer
	FEATURE  logContainer
	DOCS     logContainer
	CHORE    logContainer
	TEST     logContainer
	OTHER    logContainer
}

// printLog - loops through the collected logs to write them to string builder
func (logs logContainer) printLog(out *strings.Builder, title string) {
	if len(logs) > 0 {
		out.WriteString(fmt.Sprintf("\n\n## %s  \n", title))
		for _, item := range logs {
			out.WriteString(item + "\n")
		}
	}
}

// ToMarkdown - Generate markdown output for the collected commits
func (logs *logsByCategory) ToMarkdown() string {
	var markdownString strings.Builder

	markdownString.WriteString("# Changelog \n")

	logs.FEATURE.printLog(&markdownString, "Features")
	logs.FIX.printLog(&markdownString, "Fixes")
	logs.REFACTOR.printLog(&markdownString, "Performance")
	logs.CI.printLog(&markdownString, "CI")
	logs.DOCS.printLog(&markdownString, "Docs")
	logs.CHORE.printLog(&markdownString, "Chores")
	logs.TEST.printLog(&markdownString, "Tests")
	logs.OTHER.printLog(&markdownString, "Other Changes")

	return markdownString.String()
}
