package labeler

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

// noHiddenEnabled is a global flag to control whether hidden files should be excluded
var noHiddenEnabled bool

// SetNoHidden sets the global flag for excluding hidden files from glob matching
func SetNoHidden(enabled bool) {
	noHiddenEnabled = enabled
}

func isExtglobEnabled() bool {
	return !parser.IsEnableEnvFlag("GH_LABEL_KIT_LABELER_DISABLE_EXTGLOB")
}

// isNoHiddenEnabled checks if hidden files should be excluded from matching
func isNoHiddenEnabled() bool {
	return noHiddenEnabled
}

// isHiddenFile checks if the filename contains any hidden components (starting with .)
func isHiddenFile(filename string) bool {
	// Empty filename is not hidden
	if filename == "" {
		return false
	}
	
	// Split the path into components and check each one
	parts := strings.Split(filepath.ToSlash(filepath.Clean(filename)), "/")
	for _, part := range parts {
		if len(part) > 0 && part[0] == '.' {
			return true
		}
	}
	return false
}

func matchGlob(pattern, filename string) bool {
	// Exclude hidden files if no-hidden option is enabled
	if isNoHiddenEnabled() && isHiddenFile(filename) {
		logger.Debug("Skipping hidden file", "filename", filename)
		return false
	}

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
