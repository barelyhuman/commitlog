// SPDX-License-Identifier: MIT

package logcategory

import (
	"io"
	"strings"
)

// LogsByCategory - Type to hold logs by each's category
type LogsByCategory map[string][]string

func NewLogsByCategory() LogsByCategory {
	return make(LogsByCategory, len(Categories))
}
func (logs LogsByCategory) Add(category, message string) {
	k := strings.ToLower(category)
	logs[k] = append(logs[k], message)
}

// Categories is the Poor Man's Ordered Map of category names and titles.
var Categories = [][2]string{
	{"ci", "CI Changes"},
	{"fix", "Fixes"},
	{"refactor", "Performance Fixes"},
	{"feature", "Feature fixes"},
	{"docs", "Doc Updates"},
	{"", "Other changes"},
}

// WriteMarkdown - Generate markdown output for the collected commits
func WriteMarkdown(w io.Writer, logs LogsByCategory) error {
	ew := &errWriter{w: w}
	ew.WriteString("# Changelog  \n")

	seen := make(map[string]struct{}, len(Categories))
	var token struct{}
	var otherPrinted bool
	for _, kv := range Categories {
		vv := logs[kv[0]]
		if len(vv) != 0 {
			seen[kv[0]] = token
			ew.WriteString("\n\n## " + kv[1] + "\n")
			for _, item := range vv {
				ew.WriteString(item + "\n")
			}
			if kv[0] == "" {
				otherPrinted = true
			}
		}
	}

	for k, vv := range logs {
		if _, ok := seen[k]; ok {
			continue
		}
		if !otherPrinted {
			for _, kv := range Categories {
				if kv[0] == "" {
					ew.WriteString("\n\n## " + kv[1] + "\n")
					otherPrinted = true
					break
				}
			}
		}
		for _, item := range vv {
			ew.WriteString(item + "\n")
		}
	}

	return ew.err
}

type errWriter struct {
	err error
	w   io.Writer
}

func (ew *errWriter) Write(p []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	n, err := ew.w.Write(p)
	ew.err = err
	return n, err
}
func (ew *errWriter) WriteString(s string) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	n, err := io.WriteString(ew.w, s)
	ew.err = err
	return n, err
}
