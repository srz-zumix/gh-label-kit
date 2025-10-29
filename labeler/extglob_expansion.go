package labeler

import "strings"

// extractNegatedPatterns extracts the negated patterns from the original pattern
func extractNegatedPatterns(pattern string) ([]string, bool) {
	// Extract patterns like !(test) and convert them to positive patterns for exclusion
	var negatedPatterns []string

	// Find all negation patterns
	i := 0
	for i < len(pattern) {
		if i < len(pattern)-1 && pattern[i] == '!' && pattern[i+1] == '(' {
			// Found negation pattern
			parenEnd := findMatchingParen(pattern, i+1)
			if parenEnd == -1 {
				i++
				continue
			}

			// Extract the negated content
			negatedContent := pattern[i+2 : parenEnd]

			// Get the context before and after the negation
			before := pattern[:i]
			after := pattern[parenEnd+1:]

			// Create positive pattern for what should be excluded
			alternatives := splitExtglobAlternatives(negatedContent)
			for _, alt := range alternatives {
				alt = strings.TrimSpace(alt)
				negatedPattern := before + alt + after
				negatedPatterns = append(negatedPatterns, negatedPattern)
			}

			i = parenEnd + 1
		} else {
			i++
		}
	}

	return negatedPatterns, len(negatedPatterns) > 0
}

// replaceNegationWithWildcard replaces negation patterns with wildcards for general structure matching
func replaceNegationWithWildcard(pattern string) string {
	result := pattern

	// Find all negation patterns and replace them with *
	i := 0
	for i < len(result) {
		if i < len(result)-1 && result[i] == '!' && result[i+1] == '(' {
			// Found negation pattern
			parenEnd := findMatchingParen(result, i+1)
			if parenEnd == -1 {
				i++
				continue
			}

			// Replace the entire !(content) with *
			before := result[:i]
			after := result[parenEnd+1:]
			result = before + "*" + after

			// Continue from the position after the replacement
			i = len(before) + 1
		} else {
			i++
		}
	}

	return result
}
