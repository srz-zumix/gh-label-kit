package labeler

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/dlclark/regexp2"
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

// matchComplexGlob handles patterns that contain extglob patterns mixed with regular glob patterns
func matchComplexGlob(pattern, filename string) bool {
	// If the entire pattern is an extglob, use the existing function
	if isEntirelyExtglob(pattern) {
		return matchExtglob(pattern, filename)
	}

	// Handle negation patterns by checking exclusions first (simpler and more reliable)
	if negatedPatterns, hasNegated := extractNegatedPatterns(pattern); hasNegated {
		// Debug output
		// fmt.Printf("DEBUG: pattern=%s, filename=%s\n", pattern, filename)
		// fmt.Printf("DEBUG: negatedPatterns=%v\n", negatedPatterns)

		// First check if it matches any negated pattern (should not match if it does)
		for _, p := range negatedPatterns {
			// The negated pattern might also contain extglob, so use recursive matching
			var matched bool
			if containsExtglob(p) {
				matched = matchComplexGlob(p, filename)
			} else {
				m, err := doublestar.PathMatch(p, filename)
				matched = err == nil && m
			}
			if matched {
				return false // Matches a negated pattern, so should not match overall
			}
		}

		// If it doesn't match any negated pattern, check if it matches the general structure
		// Convert pattern like **/!(test)/**/*.go to **/*/**/*.go for structure matching
		generalPattern := replaceNegationWithWildcard(pattern)

		// Debug output
		// fmt.Printf("DEBUG: generalPattern=%s\n", generalPattern)

		// The general pattern might also contain extglob, so handle recursively
		var matched bool
		if containsExtglob(generalPattern) {
			// To avoid infinite recursion, check if it's the same as the original pattern
			if generalPattern != pattern {
				matched = matchComplexGlob(generalPattern, filename)
			} else {
				// Fall back to regex approach
				regexPattern, regexOk := convertComplexPatternToRegex2(generalPattern)
				if regexOk {
					re, err := regexp2.Compile(regexPattern, 0)
					if err == nil {
						m, err := re.MatchString(filename)
						matched = err == nil && m
					}
				}
			}
		} else {
			m, err := doublestar.PathMatch(generalPattern, filename)
			matched = err == nil && m
		}

		if matched {
			return true // Matches general structure and doesn't match negated patterns
		}
		return false
	}

	// Try regex approach for complex patterns
	regexPattern, regexOk := convertComplexPatternToRegex2(pattern)
	if regexOk {
		re, err := regexp2.Compile(regexPattern, 0)
		if err == nil {
			matched, err := re.MatchString(filename)
			if err == nil {
				return matched
			}
		}
	}

	// Fallback to regular doublestar matching
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}
