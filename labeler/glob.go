package labeler

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func isExtglobEnabled() bool {
	return !parser.IsEnableEnvFlag("GH_LABEL_KIT_LABELER_DISABLE_EXTGLOB")
}

func matchGlob(pattern, filename string) bool {
	// Check if the pattern contains any extglob patterns
	if isExtglobEnabled() && containsExtglob(pattern) {
		return matchComplexGlob(pattern, filename)
	}

	// Use regular doublestar matching for standard patterns
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}
