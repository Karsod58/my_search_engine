package search

import (
	"strings"
)

func (s *Searcher) bestSnippet(docText string, queryTerms []string, windowSize int) string {

	tokens := strings.Fields(docText)

	if len(tokens) <= windowSize {
		return docText
	}

	bestScore := 0
	bestStart := 0

	for start := 0; start <= len(tokens)-windowSize; start++ {

		end := start + windowSize
		window := tokens[start:end]

		score := scoreWindow(window, queryTerms)

		if score > bestScore {
			bestScore = score
			bestStart = start
		}
	}

	bestWindow := tokens[bestStart : bestStart+windowSize]
	return strings.Join(bestWindow, " ")
}
func scoreWindow(window []string, queryTerms []string) int {

	score := 0
	termSet := make(map[string]bool)

	for _, t := range queryTerms {
		termSet[strings.ToLower(t)] = true
	}

	for _, token := range window {
		if termSet[strings.ToLower(token)] {
			score += 2
		}
	}

	return score
}