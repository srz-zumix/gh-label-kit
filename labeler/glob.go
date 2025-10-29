package labeler

import "github.com/bmatcuk/doublestar/v4"

func matchGlob(pattern, filename string) bool {
	// Check if it's an extended glob pattern
	if isExtglob(pattern) {
		return matchExtglob(pattern, filename)
	}

	// Use regular doublestar matching for standard patterns
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}
