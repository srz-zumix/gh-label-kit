package labeler

import (
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/dlclark/regexp2"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

// ExtglobType represents the type of extended glob pattern
type ExtglobType int

const (
	ExtglobNone       ExtglobType = iota
	ExtglobNot                    // !(pattern) - match anything except pattern
	ExtglobZeroOrOne              // ?(pattern) - match zero or one occurrence of pattern
	ExtglobOneOrMore              // +(pattern) - match one or more occurrences of pattern
	ExtglobZeroOrMore             // *(pattern) - match zero or more occurrences of pattern
	ExtglobExact                  // @(pattern) - match exactly one occurrence of pattern
)

// ExtglobPattern represents a parsed extended glob pattern
type ExtglobPattern struct {
	Type     ExtglobType
	Pattern  string
	Original string
}

// parseExtglob parses an extended glob pattern with support for nested patterns
func parseExtglob(pattern string) *ExtglobPattern {
	// Try to find the first extglob pattern
	pos := findExtglobStart(pattern)
	if pos == -1 {
		return nil
	}

	// Extract the operator
	operator := pattern[pos]
	if pos+1 >= len(pattern) || pattern[pos+1] != '(' {
		return nil
	}

	// Find the matching closing parenthesis, accounting for nesting
	closePos := findMatchingParen(pattern, pos+1)
	if closePos == -1 {
		return nil
	}

	var extType ExtglobType
	switch operator {
	case '!':
		extType = ExtglobNot
	case '?':
		extType = ExtglobZeroOrOne
	case '+':
		extType = ExtglobOneOrMore
	case '*':
		extType = ExtglobZeroOrMore
	case '@':
		extType = ExtglobExact
	default:
		return nil
	}

	innerPattern := pattern[pos+2 : closePos]
	return &ExtglobPattern{
		Type:     extType,
		Pattern:  innerPattern,
		Original: pattern,
	}
}

// findExtglobStart finds the position of the first extglob operator in the pattern
func findExtglobStart(pattern string) int {
	for i := 0; i < len(pattern); i++ {
		if i+1 < len(pattern) && pattern[i+1] == '(' {
			switch pattern[i] {
			case '!', '?', '+', '*', '@':
				return i
			}
		}
	}
	return -1
}

// findMatchingParen finds the matching closing parenthesis for the opening parenthesis at openPos
func findMatchingParen(pattern string, openPos int) int {
	if openPos >= len(pattern) || pattern[openPos] != '(' {
		return -1
	}

	depth := 1
	for i := openPos + 1; i < len(pattern); i++ {
		switch pattern[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// matchComplexGlob handles patterns that contain extglob patterns mixed with regular glob patterns
func matchComplexGlob(pattern, filename string) bool {
	// If the entire pattern is an extglob, use the existing function
	if isEntirelyExtglob(pattern) {
		return matchExtglob(pattern, filename)
	}

	// Handle negation patterns by checking exclusions first (simpler and more reliable)
	if negatedPatterns, hasNegated := extractNegatedPatterns(pattern); hasNegated {
		logger.Debug("Processing complex glob with negation patterns",
			"pattern", pattern,
			"filename", filename,
			"negatedPatterns", negatedPatterns)

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

		logger.Debug("Checking general pattern structure", "generalPattern", generalPattern)

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

// matchExtglob matches a filename against an extended glob pattern with recursive support
func matchExtglob(pattern, filename string) bool {
	// Try to convert the extglob pattern to a regular expression using regexp2
	regex, ok := convertExtglobToRegex2(pattern)
	if ok {
		re, err := regexp2.Compile(regex, 0)
		if err == nil {
			matched, err := re.MatchString(filename)
			if err == nil {
				return matched
			}
		}
		// Log debug information when regex compilation or matching fails
		logger.Debug("Extglob regex matching failed", "pattern", pattern, "regex", regex, "filename", filename)
	}

	// Fallback to the original implementation
	extPattern := parseExtglob(pattern)
	if extPattern == nil {
		// Not an extended glob, use regular doublestar matching
		return matchGlobDoublestar(pattern, filename)
	}

	// Special handling for negation patterns with remaining parts
	if extPattern.Type == ExtglobNot && len(extPattern.Original) > len(extPattern.Pattern)+4 {
		// Pattern like !(test)/* - handle the remaining part
		return matchExtglobNotWithRemainder(pattern, filename)
	}

	switch extPattern.Type {
	case ExtglobNot:
		return matchExtglobNot(extPattern.Pattern, filename)
	case ExtglobZeroOrOne:
		return matchExtglobZeroOrOne(extPattern.Pattern, filename)
	case ExtglobOneOrMore:
		return matchExtglobOneOrMore(extPattern.Pattern, filename)
	case ExtglobZeroOrMore:
		return matchExtglobZeroOrMore(extPattern.Pattern, filename)
	case ExtglobExact:
		return matchExtglobExact(extPattern.Pattern, filename)
	default:
		return false
	}
}

// matchPattern matches a pattern that may contain nested extglob patterns
func matchPattern(pattern, filename string) bool {
	// Check if the pattern contains extglob syntax
	if isExtglob(pattern) {
		return matchExtglob(pattern, filename)
	}
	// Use regular doublestar matching for standard patterns
	return matchGlobDoublestar(pattern, filename)
}

// splitExtglobAlternatives splits a pattern by | while respecting nested parentheses
func splitExtglobAlternatives(pattern string) []string {
	if !strings.Contains(pattern, "|") {
		return []string{pattern}
	}

	var alternatives []string
	var current strings.Builder
	depth := 0

	for _, r := range pattern {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case '|':
			if depth == 0 {
				// We're at the top level, this is a separator
				alternatives = append(alternatives, current.String())
				current.Reset()
			} else {
				// We're inside nested parentheses, include the |
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add the last alternative
	if current.Len() > 0 {
		alternatives = append(alternatives, current.String())
	}

	return alternatives
}

// convertExtglobToRegex2 converts an extglob pattern to a regexp2-compatible regular expression
func convertExtglobToRegex2(pattern string) (string, bool) {
	if !isExtglob(pattern) {
		return "", false
	}

	// Process extglob patterns manually to avoid regex complexity
	converted := processExtglobManually(pattern)

	// Add anchors if not already present (for non-negation patterns)
	if !strings.HasPrefix(converted, "^(?!") {
		// Add anchors if not already present
		if !strings.HasPrefix(converted, "^") {
			converted = "^" + converted
		}
		if !strings.HasSuffix(converted, "$") {
			converted = converted + "$"
		}
	}

	return converted, true
}

// convertGlobToRegex converts basic glob patterns to regex with proper path handling
func convertGlobToRegex(pattern string) string {
	var result strings.Builder

	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		switch char {
		case '*':
			// Check for ** pattern
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				// ** matches any number of directories (including zero)
				// Check if there's a / after **
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					// **/ pattern - matches any path with / or nothing
					result.WriteString("(?:[^/]*/)*")
					i += 2 // Skip the next * and /
				} else {
					// ** at end or followed by non-/ - matches anything
					result.WriteString(".*")
					i++ // Skip the next *
				}
			} else {
				// Single * matches any characters except path separator
				result.WriteString("[^/]*")
			}
		case '?':
			// ? matches any single character except path separator
			result.WriteString("[^/]")
		case '.', '^', '$', '{', '}', '[', ']', '\\':
			// Escape regex special characters (but not +, |, or parentheses)
			result.WriteString("\\")
			result.WriteByte(char)
		default:
			result.WriteByte(char)
		}
	}

	return result.String()
}

// matchGlobDoublestar is the original glob matching using doublestar
func matchGlobDoublestar(pattern, filename string) bool {
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}

// matchExtglobNot implements !(pattern) - match anything except pattern
func matchExtglobNot(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := splitExtglobAlternatives(pattern)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchPattern(p, filename) {
			return false
		}
	}
	return true
}

// matchExtglobNotWithRemainder handles negation patterns with remaining parts like !(test)/*
func matchExtglobNotWithRemainder(pattern, filename string) bool {
	// Parse the pattern to extract negation part and remainder
	extPattern := parseExtglob(pattern)
	if extPattern == nil || extPattern.Type != ExtglobNot {
		return false
	}

	// Find the remainder part after the closing parenthesis
	pos := findExtglobStart(pattern)
	closePos := findMatchingParen(pattern, pos+1)
	if closePos == -1 || closePos+1 >= len(pattern) {
		return matchExtglobNot(extPattern.Pattern, filename)
	}

	remainder := pattern[closePos+1:]

	// Check if filename matches the overall structure (with remainder)
	// Create a generic pattern that matches the same structure
	genericPattern := "*" + remainder
	structureMatches := matchGlobDoublestar(genericPattern, filename)

	logger.Debug("Extglob pattern matching with remainder",
		"pattern", pattern,
		"filename", filename,
		"remainder", remainder,
		"genericPattern", genericPattern,
		"structureMatches", structureMatches)

	if !structureMatches {
		return false // Doesn't match the required structure
	}

	// Check if filename matches any of the negated patterns + remainder
	patterns := splitExtglobAlternatives(extPattern.Pattern)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		negatedPattern := p + remainder
		if matchGlobDoublestar(negatedPattern, filename) {
			logger.Debug("Filename matches negated pattern", "negatedPattern", negatedPattern)
			return false // Matches a negated pattern
		}
	}

	return true // Matches structure but doesn't match negated patterns
}

// matchExtglobZeroOrOne implements ?(pattern) - match zero or one occurrence
func matchExtglobZeroOrOne(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := splitExtglobAlternatives(pattern)

	// Check if filename matches any pattern (one occurrence)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchPattern(p, filename) {
			return true
		}
	}

	// Check if filename is empty or matches the "zero occurrence" case
	// For simplicity, we consider zero occurrence as empty string or base pattern without the extglob part
	return filename == "" || matchPattern("", filename)
}

// matchExtglobOneOrMore implements +(pattern) - match one or more occurrences
func matchExtglobOneOrMore(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := splitExtglobAlternatives(pattern)

	// At least one pattern must match
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchPattern(p, filename) {
			return true
		}
	}

	return false
}

// matchExtglobZeroOrMore implements *(pattern) - match zero or more occurrences
func matchExtglobZeroOrMore(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := splitExtglobAlternatives(pattern)

	// Check if filename matches any pattern (one or more occurrences)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchPattern(p, filename) {
			return true
		}
	}

	// Also match zero occurrences (empty string or base pattern)
	return true // Zero or more means it can always match
}

// matchExtglobExact implements @(pattern) - match exactly one occurrence
func matchExtglobExact(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := splitExtglobAlternatives(pattern)

	// Check if filename matches exactly one of the patterns
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchPattern(p, filename) {
			return true
		}
	}

	return false
}

// isExtglob checks if a pattern contains extended glob syntax
func isExtglob(pattern string) bool {
	// Check for any extglob pattern anywhere in the string
	extglobPattern := regexp.MustCompile(`[!?+*@]\([^)]*\)`)
	return extglobPattern.MatchString(pattern)
}

// isEntirelyExtglob checks if the entire pattern is a single extglob pattern
func isEntirelyExtglob(pattern string) bool {
	// Check if pattern starts with extglob and covers the entire pattern
	if len(pattern) < 3 {
		return false
	}

	// Check if it starts with extglob operator
	if !strings.ContainsRune("!?+*@", rune(pattern[0])) || pattern[1] != '(' {
		return false
	}

	// Find the matching closing parenthesis
	parenEnd := findMatchingParen(pattern, 1)
	if parenEnd == -1 {
		return false
	}

	// Check if the pattern ends at the closing parenthesis (possibly with some trailing simple patterns)
	remaining := pattern[parenEnd+1:]

	// If remaining contains another extglob, this is not entirely a single extglob
	if containsExtglob(remaining) {
		return false
	}

	// Allow simple trailing patterns like /* or /**
	if remaining == "" || remaining == "/*" || remaining == "/**" || strings.HasPrefix(remaining, "/*") || strings.HasPrefix(remaining, "/**") {
		return true
	}

	return false
}

// String returns a string representation of the extglob type
func (t ExtglobType) String() string {
	switch t {
	case ExtglobNot:
		return "not"
	case ExtglobZeroOrOne:
		return "zero-or-one"
	case ExtglobOneOrMore:
		return "one-or-more"
	case ExtglobZeroOrMore:
		return "zero-or-more"
	case ExtglobExact:
		return "exact"
	default:
		return "none"
	}
}

// processExtglobManually manually processes extglob patterns to avoid regex complexity
func processExtglobManually(pattern string) string {
	result := ""
	i := 0

	for i < len(pattern) {
		// Look for extglob patterns
		if i < len(pattern)-1 && strings.ContainsRune("!?+*@", rune(pattern[i])) && pattern[i+1] == '(' {
			operator := pattern[i]

			// Find matching closing parenthesis
			parenStart := i + 1
			parenEnd := findMatchingParen(pattern, parenStart)

			if parenEnd == -1 {
				// No matching paren, treat as literal
				result += convertGlobToRegex(string(pattern[i]))
				i++
				continue
			}

			// Extract content
			content := pattern[i+2 : parenEnd]

			// Process content
			alternatives := splitExtglobAlternatives(content)
			regexAlternatives := make([]string, len(alternatives))
			for j, alt := range alternatives {
				alt = strings.TrimSpace(alt)
				// Recursively process nested extglob in alternatives
				if containsExtglob(alt) {
					regexAlternatives[j] = processExtglobManually(alt)
				} else {
					regexAlternatives[j] = convertGlobToRegex(alt)
				}
			}
			regexContent := strings.Join(regexAlternatives, "|")

			// Apply operator
			switch operator {
			case '!':
				// Special handling for negation
				remainingPattern := pattern[parenEnd+1:]
				if remainingPattern != "" {
					// For patterns like !(test)/* or !(test)/**
					// We need to match the structure but exclude specific patterns
					convertedRemainder := convertGlobToRegex(remainingPattern)

					// First part: ensure it doesn't match the negated pattern + remainder
					negatedFull := "^(?!(" + regexContent + ")" + convertedRemainder + "$)"

					// Second part: match the general structure
					generalPattern := convertGlobToRegex("*" + remainingPattern)

					// Combine: negative lookahead + positive match for structure
					return negatedFull + "(" + generalPattern + ")$"
				} else {
					return "^(?!(" + regexContent + ")$).*$"
				}
			case '?':
				result += "(" + regexContent + ")?"
			case '+':
				result += "(" + regexContent + ")+"
			case '*':
				result += "(" + regexContent + ")*"
			case '@':
				result += "(" + regexContent + ")"
			}

			i = parenEnd + 1
		} else {
			// Regular character - handle with proper path glob conversion
			char := pattern[i]
			switch char {
			case '.', '^', '$', '{', '}', '[', ']', '\\':
				result += "\\" + string(char)
			case '*':
				// Check for ** pattern
				if i+1 < len(pattern) && pattern[i+1] == '*' {
					// ** matches any number of directories (including zero)
					// Check if there's a / after **
					if i+2 < len(pattern) && pattern[i+2] == '/' {
						// **/ pattern - matches any path with / or nothing
						result += "(?:[^/]*/)*"
						i += 2 // Skip the next * and /
					} else {
						// ** at end or followed by non-/ - matches anything
						result += ".*"
						i++ // Skip the next *
					}
				} else {
					// Single * matches any characters except path separator
					result += "[^/]*"
				}
			case '?':
				// ? matches any single character except path separator
				result += "[^/]"
			default:
				result += string(char)
			}
			i++
		}
	}

	return result
}

// containsExtglob checks if a pattern contains extglob without regex
func containsExtglob(pattern string) bool {
	for i := 0; i < len(pattern)-1; i++ {
		if strings.ContainsRune("!?+*@", rune(pattern[i])) && pattern[i+1] == '(' {
			return true
		}
	}
	return false
}

// convertComplexPatternToRegex2 converts complex patterns with extglob to regexp2
func convertComplexPatternToRegex2(pattern string) (string, bool) {
	// Process the pattern manually to handle mixed extglob and regular patterns
	result := processComplexPattern(pattern)
	if result == "" {
		return "", false
	}

	// Add anchors
	if !strings.HasPrefix(result, "^") {
		result = "^" + result
	}
	if !strings.HasSuffix(result, "$") {
		result = result + "$"
	}

	return result, true
}

// processComplexPattern processes patterns with mixed extglob and regular glob
func processComplexPattern(pattern string) string {
	result := ""
	i := 0

	for i < len(pattern) {
		// Look for extglob patterns
		if i < len(pattern)-1 && strings.ContainsRune("!?+*@", rune(pattern[i])) && pattern[i+1] == '(' {
			// Found extglob, process it
			operator := pattern[i]
			parenStart := i + 1
			parenEnd := findMatchingParen(pattern, parenStart)

			if parenEnd == -1 {
				// No matching paren, treat as literal
				result += convertGlobToRegex(string(pattern[i]))
				i++
				continue
			}

			// Extract content
			content := pattern[i+2 : parenEnd]

			// Process content
			alternatives := splitExtglobAlternatives(content)
			regexAlternatives := make([]string, len(alternatives))
			for j, alt := range alternatives {
				alt = strings.TrimSpace(alt)
				// Recursively process nested patterns
				if containsExtglob(alt) {
					regexAlternatives[j] = processComplexPattern(alt)
				} else {
					regexAlternatives[j] = convertGlobToRegex(alt)
				}
			}
			regexContent := strings.Join(regexAlternatives, "|")

			// Apply operator - special handling for negation with mixed patterns
			switch operator {
			case '!':
				// For negation in complex patterns, we need to be more careful
				// Get what comes after this extglob
				remainingPattern := pattern[parenEnd+1:]
				if remainingPattern != "" {
					// Convert remainder to regex - handle nested extglob recursively
					var convertedRemainder string
					if containsExtglob(remainingPattern) {
						convertedRemainder = processComplexPattern(remainingPattern)
					} else {
						convertedRemainder = convertGlobToRegex(remainingPattern)
					}
					// Use negative lookahead for the entire negated pattern + remainder
					// but also provide a positive pattern to match against
					result += "(?!(" + regexContent + ")" + convertedRemainder + "$)[^/]*" + convertedRemainder
				} else {
					result += "(?!(" + regexContent + ")$)"
				}
			case '?':
				result += "(" + regexContent + ")?"
			case '+':
				result += "(" + regexContent + ")+"
			case '*':
				result += "(" + regexContent + ")*"
			case '@':
				result += "(" + regexContent + ")"
			}

			i = parenEnd + 1
		} else {
			// Regular character or glob pattern
			char := pattern[i]
			switch char {
			case '.', '^', '$', '{', '}', '[', ']', '\\':
				result += "\\" + string(char)
			case '*':
				// Check for ** pattern
				if i+1 < len(pattern) && pattern[i+1] == '*' {
					// ** matches any number of directories
					if i+2 < len(pattern) && pattern[i+2] == '/' {
						// **/ pattern
						result += "(?:[^/]*/)*"
						i += 2
					} else {
						// ** at end or followed by non-/
						result += ".*"
						i++
					}
				} else {
					// Single * matches any characters except path separator
					result += "[^/]*"
				}
			case '?':
				// ? matches any single character except path separator
				result += "[^/]"
			default:
				result += string(char)
			}
			i++
		}
	}

	return result
}

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
