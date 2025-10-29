package labeler

import (
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
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

// extglobRegex matches extended glob patterns
var extglobRegex = regexp.MustCompile(`^([!?+*@])\(([^)]*)\)(.*)$`)

// parseExtglob parses an extended glob pattern
func parseExtglob(pattern string) *ExtglobPattern {
	matches := extglobRegex.FindStringSubmatch(pattern)
	if matches == nil {
		return nil
	}

	var extType ExtglobType
	switch matches[1] {
	case "!":
		extType = ExtglobNot
	case "?":
		extType = ExtglobZeroOrOne
	case "+":
		extType = ExtglobOneOrMore
	case "*":
		extType = ExtglobZeroOrMore
	case "@":
		extType = ExtglobExact
	default:
		return nil
	}

	return &ExtglobPattern{
		Type:     extType,
		Pattern:  matches[2],
		Original: pattern,
	}
}

// matchExtglob matches a filename against an extended glob pattern
func matchExtglob(pattern, filename string) bool {
	extPattern := parseExtglob(pattern)
	if extPattern == nil {
		// Not an extended glob, use regular doublestar matching
		return matchGlobDoublestar(pattern, filename)
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

// matchGlobDoublestar is the original glob matching using doublestar
func matchGlobDoublestar(pattern, filename string) bool {
	matched, err := doublestar.PathMatch(pattern, filename)
	return err == nil && matched
}

// matchExtglobNot implements !(pattern) - match anything except pattern
func matchExtglobNot(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := strings.Split(pattern, "|")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchGlobDoublestar(p, filename) {
			return false
		}
	}
	return true
}

// matchExtglobZeroOrOne implements ?(pattern) - match zero or one occurrence
func matchExtglobZeroOrOne(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := strings.Split(pattern, "|")

	// Check if filename matches any pattern (one occurrence)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchGlobDoublestar(p, filename) {
			return true
		}
	}

	// Check if filename is empty or matches the "zero occurrence" case
	// For simplicity, we consider zero occurrence as empty string or base pattern without the extglob part
	return filename == "" || matchGlobDoublestar("", filename)
}

// matchExtglobOneOrMore implements +(pattern) - match one or more occurrences
func matchExtglobOneOrMore(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := strings.Split(pattern, "|")

	// At least one pattern must match
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchGlobDoublestar(p, filename) {
			return true
		}
	}

	return false
}

// matchExtglobZeroOrMore implements *(pattern) - match zero or more occurrences
func matchExtglobZeroOrMore(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := strings.Split(pattern, "|")

	// Check if filename matches any pattern (one or more occurrences)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchGlobDoublestar(p, filename) {
			return true
		}
	}

	// Also match zero occurrences (empty string or base pattern)
	return true // Zero or more means it can always match
}

// matchExtglobExact implements @(pattern) - match exactly one occurrence
func matchExtglobExact(pattern, filename string) bool {
	// Handle multiple patterns separated by |
	patterns := strings.Split(pattern, "|")

	// Check if filename matches exactly one of the patterns
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if matchGlobDoublestar(p, filename) {
			return true
		}
	}

	return false
}

// isExtglob checks if a pattern contains extended glob syntax
func isExtglob(pattern string) bool {
	return extglobRegex.MatchString(pattern)
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
