package labeler

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func isExtglobEnabled() bool {
	return !parser.IsEnableEnvFlag("GH_LABEL_KIT_LABELER_DISABLE_EXTGLOB")
}

func matchGlob(pattern, filename string) bool {
	// Check if the pattern contains any extglob patterns
	if isExtglobEnabled() && containsExtglob(pattern) {
		result := matchComplexGlob(pattern, filename)
		logger.Debug("Extglob pattern match", "pattern", pattern, "filename", filename, "matched", result)
		return result
	}

	// Check for negation pattern (! at the beginning) only when extglob is enabled
	negate := false
	if len(pattern) > 0 && pattern[0] == '!' {
		negate = true
		pattern = pattern[1:]
	}

	// Use regular doublestar matching for standard patterns
	matched, err := doublestar.PathMatch(pattern, filename)
	result := err == nil && matched

	// Invert the result if negation is enabled
	if negate {
		result = !result
	}
	logger.Debug("Glob pattern match", "pattern", pattern, "filename", filename, "negate", negate, "matched", result)
	return result
}
