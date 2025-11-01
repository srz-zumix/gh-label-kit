package labeler

import (
	"github.com/bmatcuk/doublestar/v4"
)

func matchGlob(pattern, filename string) bool {
	// Check if the pattern contains any extglob patterns
	if containsExtglob(pattern) {
		return matchComplexGlob(pattern, filename)
	}

	// Use regular doublestar matching for standard patterns
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}
