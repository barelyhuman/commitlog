// SPDX-License-Identifier: MIT

package logcategory

import (
	"io"
)

// LogsByCategory - Type to hold logs by each's category
type LogsByCategory struct {
	CI       []string
	Fix      []string
	Refactor []string
	Feature  []string
	Docs     []string
	Other    []string
}

// WriteMarkdown - Generate markdown output for the collected commits
func WriteMarkdown(w io.Writer, logs *LogsByCategory) error {
	ew := &errWriter{w: w}
	ew.WriteString("# Changelog  \n")

	if len(logs.CI) != 0 {
		ew.WriteString("\n\n## CI Changes  \n")
		for _, item := range logs.CI {
			ew.WriteString(item + "\n")
		}
	}

	if len(logs.Fix) != 0 {
		ew.WriteString("\n\n## Fixes  \n")
		for _, item := range logs.Fix {
			ew.WriteString(item + "\n")
		}
	}

	if len(logs.Refactor) != 0 {
		ew.WriteString("\n\n## Performance Fixes  \n")
		for _, item := range logs.Refactor {
			ew.WriteString(item + "\n")
		}
	}

	if len(logs.Feature) != 0 {
		ew.WriteString("\n\n## Feature Fixes  \n")
		for _, item := range logs.Feature {
			ew.WriteString(item + "\n")
		}
	}

	if len(logs.Docs) != 0 {
		ew.WriteString("\n\n## Doc Updates  \n")
		for _, item := range logs.Docs {
			ew.WriteString(item + "\n")
		}
	}

	if len(logs.Other) != 0 {
		ew.WriteString("\n\n## Other Changes  \n")
		for _, item := range logs.Other {
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
