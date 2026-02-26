package search

import "strings"

type AutoComplete struct {
	terms []string
}

func NewAutoComplete(vocab []string) *AutoComplete {
	return &AutoComplete{terms: vocab}
}

func (a *AutoComplete) Suggest(prefix string, limit int) []string {
	if prefix == "" {
		return nil
	}

	prefix = strings.ToLower(prefix)

	out := []string{}

	for _, t := range a.terms {
		if strings.HasPrefix(t, prefix) {
			out = append(out, t)
			if len(out) >= limit {
				break
			}
		}
	}

	return out
}