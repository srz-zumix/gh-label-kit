package labeler

import "github.com/bmatcuk/doublestar/v4"

func matchGlob(pattern, filename string) bool {
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}
