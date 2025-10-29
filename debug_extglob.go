package main

import (
	"fmt"

	"github.com/srz-zumix/gh-label-kit/labeler"
)

func main() {
	pattern := "!(node_modules)/**/*.@(js|ts)"
	filename := "src/app.js"

	// Test extractNegatedPatterns
	negatedPatterns, hasNegated := labeler.ExtractNegatedPatterns(pattern)
	fmt.Printf("Pattern: %s\n", pattern)
	fmt.Printf("HasNegated: %v\n", hasNegated)
	fmt.Printf("NegatedPatterns: %v\n", negatedPatterns)

	// Test replaceNegationWithWildcard
	generalPattern := labeler.ReplaceNegationWithWildcard(pattern)
	fmt.Printf("GeneralPattern: %s\n", generalPattern)

	// Test matchGlob
	result := labeler.MatchGlob(pattern, filename)
	fmt.Printf("MatchResult: %v\n", result)
}
