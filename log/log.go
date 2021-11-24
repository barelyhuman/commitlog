// SPDX-License-Identifier: MIT

package commitlog

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// SupportedKeys - keys that are supported by the package
const SupportedKeys = "ci|refactor|docs|fix|feat|test|chore|other"

// ErrMessage - simple interface around error with a custom message
type ErrMessage struct {
	Message string
	Err     error
}

type commitTypeInclusions [][]byte

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
			out.WriteString(item + "  \n")
		}
	}
}

// ToMarkdown - Generate markdown output for the collected commits
func (logs *logsByCategory) ToMarkdown(skipped bool) string {
	var markdownString strings.Builder

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

/*
	TODO:
	- [] if the current start is also a tag then get data till prev tag
	- [] add in option to include the description, if the commit has a description
*/

// CommitLog - Generate commit log
func CommitLog(path string, startCommitString string, endCommitString string, inclusionFlags string, skipClassification bool) (string, ErrMessage) {
	currentRepository := OpenRepository(path)
	baseCommitReference, err := currentRepository.Head()
	var startHash, endHash *object.Commit
	var cIter object.CommitIter

	if err != nil {
		return "", ErrMessage{"Unable to get repository HEAD:", err}
	}

	startHash = GetCommitFromString(startCommitString, currentRepository)
	endHash = GetCommitFromString(endCommitString, currentRepository)

	if startHash != nil {
		cIter, err = currentRepository.Log(&git.LogOptions{From: startHash.Hash})
	} else {
		cIter, err = currentRepository.Log(&git.LogOptions{From: baseCommitReference.Hash()})
	}

	if err != nil {
		return "", ErrMessage{"Unable to get repository commits:", err}
	}

	var commits []*object.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	if err != nil {
		return "", ErrMessage{"Error getting commits : ", err}
	}

	var inclusions commitTypeInclusions
	inclusions = append(inclusions, bytes.SplitN([]byte(inclusionFlags), []byte("|"), -1)...)
	inclusions = append(inclusions, bytes.SplitN([]byte(inclusionFlags), []byte(","), -1)...)
	logContainer := logsByCategory{}

	logContainer.Setup()

	logContainer.CI.include = inclusions.checkInclusion("ci")
	logContainer.FIX.include = inclusions.checkInclusion("fix")
	logContainer.REFACTOR.include = inclusions.checkInclusion("refactor")
	logContainer.FEATURE.include = inclusions.checkInclusion("feat")
	logContainer.DOCS.include = inclusions.checkInclusion("docs")
	logContainer.CHORE.include = inclusions.checkInclusion("chore")
	logContainer.TEST.include = inclusions.checkInclusion("test")
	logContainer.OTHER.include = inclusions.checkInclusion("other")

	for _, c := range commits {
		key, scopedKey := findKeyInCommit(SupportedKeys, c.Message)
		key = strings.SplitN(strings.TrimSpace(key), ":", 2)[0]
		normalizedHash := c.Hash.String() + " - " + normalizeCommit(c.Message, scopedKey)

		logContainer.AddCommit(key, normalizedHash, skipClassification)

		if endHash == nil && isCommitToNearestTag(currentRepository, c) {
			break
		} else if endHash != nil && c.Hash == endHash.Hash {
			break
		}
	}

	return logContainer.ToMarkdown(skipClassification), ErrMessage{}
}

func (inclusions *commitTypeInclusions) checkInclusion(flagToCheck string) bool {
	if inclusions != nil {
		for _, flag := range *inclusions {
			if string(flag) == flagToCheck {
				return true
			}
		}
	}
	return false
}

func findKeyInCommit(key string, commitMessage string) (string, string) {
	re := regexp.MustCompile(`^(` + key + `)[:]|^((` + key + `)\((\w+[, /\\]*)+\)[:])`)
	keyMatches := re.FindAllStringSubmatch(commitMessage, -1)
	if len(keyMatches) == 0 {
		return "", ""
	}

	scopedKey := keyMatches[0][0]
	scopelessKey := keyMatches[0][3]

	if scopelessKey == "" {
		scopelessKey = keyMatches[0][1]
	}

	return scopelessKey, scopedKey
}
