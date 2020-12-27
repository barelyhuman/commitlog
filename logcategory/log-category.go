package logcategory

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
