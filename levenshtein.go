package dictionary

import "math"

func levenshteinRecursive(str1, str2 string, m, n int) int {
	// Base case: str1 is empty
	if m == 0 {
		return n
	}

	// Base case: str2 is empty
	if n == 0 {
		return m
	}

	// If the last characters of both
	// strings are the same
	if str1[m-1] == str2[n-1] {
		return levenshteinRecursive(str1, str2, m-1, n-1)
	}

	// Calculate the minimum of three possible
	// operations (insert, remove, replace)
	var minOfInsertRemove = math.Min(
		// Insert
		float64(levenshteinRecursive(str1, str2, m, n-1)),
		// Remove
		float64(levenshteinRecursive(str1, str2, m-1, n)))

	return 1 + int(math.Min(
		minOfInsertRemove,
		// Replace
		float64(levenshteinRecursive(str1, str2, m-1, n-1))),
	)
}
