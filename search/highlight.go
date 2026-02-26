package search

import "strings"

const (
	Bold  = "\033[1m"
	Cyan  = "\033[36m"
	Reset = "\033[0m"
)

func Highlight(text string, terms []string) string {
	lower := strings.ToLower(text)

	for _, t := range terms {
		if t == "" {
			continue
		}

		lowerTerm := strings.ToLower(t)

		if strings.Contains(lower, lowerTerm) {
			text = strings.ReplaceAll(
				text,
				t,
				Bold+Cyan+t+Reset,
			)
		}
	}

	return text
}