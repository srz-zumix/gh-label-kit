// Package labeler provides extended glob (extglob) pattern matching functionality
// that follows bash shell extglob behavior.
//
// Key differences from standard glob patterns:
//   - Inside extglob expressions (e.g., !(pattern), @(pattern|pattern2)), wildcards
//     (* and ?) can match across path separators (/), following bash shell behavior
//   - Outside extglob expressions, standard glob rules apply where * does not cross
//     directory boundaries
//
// Examples:
//   - !(test) matches "dir/file" (wildcard crosses path separator)
//   - !(*test*) excludes any path containing "test", including "src/test/file.go"
//   - **/*.go uses standard glob where * only matches within a directory level
package labeler

import (
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/dlclark/regexp2"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

// =============================================================================
// TYPE DEFINITIONS
// =============================================================================

// ExtglobType represents the type of extended glob pattern supported by bash.
// Each type defines different matching semantics for the enclosed pattern.
type ExtglobType int

const (
	// ExtglobNone indicates no extended glob pattern
	ExtglobNone ExtglobType = iota
	// ExtglobNot represents !(pattern) - matches anything except the specified pattern(s)
	ExtglobNot
	// ExtglobZeroOrOne represents ?(pattern) - matches zero or one occurrence of pattern(s)
	ExtglobZeroOrOne
	// ExtglobOneOrMore represents +(pattern) - matches one or more occurrences of pattern(s)
	ExtglobOneOrMore
	// ExtglobZeroOrMore represents *(pattern) - matches zero or more occurrences of pattern(s)
	ExtglobZeroOrMore
	// ExtglobExact represents @(pattern) - matches exactly one of the specified pattern(s)
	ExtglobExact
)

// String returns a string representation of the ExtglobType for debugging
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

// ExtglobPattern represents a parsed extended glob pattern with its components
type ExtglobPattern struct {
	Type     ExtglobType // The type of extglob pattern
	Pattern  string      // The inner pattern content (without the operator and parentheses)
	Original string      // The original full pattern string
}

// =============================================================================
// PATTERN PARSING UTILITIES
// =============================================================================

// parseExtglob parses an extended glob pattern with support for nested patterns.
// Returns nil if the pattern is not a valid extglob pattern.
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

// findExtglobStart finds the position of the first extglob operator in the pattern.
// Returns -1 if no extglob operator is found.
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

// findMatchingParen finds the matching closing parenthesis for the opening parenthesis at openPos.
// Properly handles nested parentheses. Returns -1 if no matching parenthesis is found.
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

// splitExtglobAlternatives splits a pattern by | while respecting nested parentheses.
// This is used to handle patterns like "js|ts|jsx" within extglob expressions.
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

// isExtglob checks if a pattern contains extended glob syntax anywhere in the string
func isExtglob(pattern string) bool {
	// Check for any extglob pattern anywhere in the string
	extglobPattern := regexp.MustCompile(`[!?+*@]\([^)]*\)`)
	return extglobPattern.MatchString(pattern)
}

// containsExtglob checks if a pattern contains extglob without using regex for performance
func containsExtglob(pattern string) bool {
	for i := 0; i < len(pattern)-1; i++ {
		if strings.ContainsRune("!?+*@", rune(pattern[i])) && pattern[i+1] == '(' {
			return true
		}
	}
	return false
}

// isEntirelyExtglob checks if the entire pattern is a single extglob pattern
// (possibly with some trailing simple glob patterns like /* or /**)
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

// =============================================================================
// MAIN MATCHING FUNCTIONS
// =============================================================================

// matchComplexGlob handles patterns that contain extglob patterns mixed with regular glob patterns.
// This is the main entry point for complex pattern matching.
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

// convertGlobToRegex converts basic glob patterns to regex with proper path handling.
// This is used for patterns OUTSIDE of extglob expressions where * should not cross path separators.
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
				// Single * matches any characters except path separator (standard glob behavior)
				result.WriteString("[^/]*")
			}
		case '?':
			// ? matches any single character except path separator (standard glob behavior)
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

// convertExtglobContentToRegex converts glob patterns within extglob expressions to regex.
// In shell extglob, wildcards inside extglob patterns CAN cross path separators.
// This matches bash's extglob behavior where !(test) can match paths like "dir/file".
func convertExtglobContentToRegex(pattern string) string {
	var result strings.Builder

	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		switch char {
		case '*':
			// In extglob context, * can cross path separators (shell behavior)
			result.WriteString(".*")
		case '?':
			// In extglob context, ? matches any single character including / (shell behavior)
			result.WriteString(".")
		case '.', '^', '$', '{', '}', '[', ']', '\\':
			// Escape regex special characters
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

// =============================================================================
// EXTGLOB TYPE-SPECIFIC MATCHING FUNCTIONS
// =============================================================================

// matchExtglobZeroOrOne implements ?(pattern) - match zero or one occurrence of the pattern.
// Returns true if filename matches any of the alternatives or represents zero occurrence.
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

// matchExtglobOneOrMore implements +(pattern) - match one or more occurrences of the pattern.
// At least one pattern alternative must match for this to return true.
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

// matchExtglobZeroOrMore implements *(pattern) - match zero or more occurrences of the pattern.
// Always matches since zero occurrences is valid, unless a specific pattern matches.
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

// matchExtglobExact implements @(pattern) - match exactly one occurrence of the pattern.
// Returns true if filename matches exactly one of the pattern alternatives.
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
					// Use extglob-specific regex conversion for content inside extglob
					// This allows * and ? to cross path separators (shell behavior)
					regexAlternatives[j] = convertExtglobContentToRegex(alt)
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
					// Shell behavior: !(test)/* should match both "main.go" and "src/main.go"
					// The negation part can match empty string, so we need to handle root files

					// Convert remainder using normal glob rules (not extglob rules)
					convertedRemainder := convertGlobToRegex(remainingPattern)

					// First negative lookahead: exclude exact matches of negated pattern + remainder
					negatedFull := "^(?!(" + regexContent + ")" + convertedRemainder + "$)"

					var generalPattern string
					// For shell compatibility, we need to allow matching from the start
					// !(test)/* should match "main.go" because !(test) can match empty string
					if strings.HasPrefix(remainingPattern, "/") {
						// Handle patterns like !(test)/* or !(test)/**
						// This should match files that start with the remainder pattern
						// but also files that have some prefix + remainder
						remainderWithoutSlash := remainingPattern[1:] // Remove leading /
						convertedWithoutSlash := convertGlobToRegex(remainderWithoutSlash)

						// Match either:
						// 1. Files that match the remainder directly (like main.go for /*)
						// 2. Files that have some path + remainder (like src/main.go for /*)
						generalPattern = "(" + convertedWithoutSlash + "|.*" + convertedRemainder + ")"
					} else {
						// For patterns without leading slash
						generalPattern = convertGlobToRegex("*" + remainingPattern)
					}

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
					// Use extglob-specific regex conversion for content inside extglob
					regexAlternatives[j] = convertExtglobContentToRegex(alt)
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

// =============================================================================
// PATTERN EXPANSION AND HELPER FUNCTIONS
// =============================================================================

// extractNegatedPatterns extracts the negated patterns from the original pattern.
// Used for complex negation handling in patterns like **/!(test)/**/*.go
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

// replaceNegationWithWildcard replaces negation patterns with wildcards for general structure matching.
// Converts patterns like **/!(test)/**/*.go to **/*/**/*.go for structural validation.
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
