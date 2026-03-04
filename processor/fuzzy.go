package processor

import "strings"

// LevenshteinDistance calculates the edit distance between two strings
func LevenshteinDistance(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Create matrix
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// FindClosest finds vocabulary words within maxDist edit distance
func FindClosest(word string, vocabulary []string, maxDist int) []string {
	var matches []string
	word = strings.ToLower(word)

	for _, vocabWord := range vocabulary {
		// Optimization: skip if length difference is too large
		lenDiff := len(word) - len(vocabWord)
		if lenDiff < 0 {
			lenDiff = -lenDiff
		}
		if lenDiff > maxDist {
			continue
		}

		dist := LevenshteinDistance(word, vocabWord)
		if dist <= maxDist {
			matches = append(matches, vocabWord)
		}
	}

	return matches
}
