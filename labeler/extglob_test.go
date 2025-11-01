package labeler

import (
	"testing"
)

func TestParseExtglob(t *testing.T) {
	cases := []struct {
		pattern  string
		expected *ExtglobPattern
		wantNil  bool
	}{
		{
			pattern: "!(*.md)",
			expected: &ExtglobPattern{
				Type:     ExtglobNot,
				Pattern:  "*.md",
				Original: "!(*.md)",
			},
		},
		{
			pattern: "?(*.js|*.ts)",
			expected: &ExtglobPattern{
				Type:     ExtglobZeroOrOne,
				Pattern:  "*.js|*.ts",
				Original: "?(*.js|*.ts)",
			},
		},
		{
			pattern: "+(test|spec)",
			expected: &ExtglobPattern{
				Type:     ExtglobOneOrMore,
				Pattern:  "test|spec",
				Original: "+(test|spec)",
			},
		},
		{
			pattern: "*(src|lib)",
			expected: &ExtglobPattern{
				Type:     ExtglobZeroOrMore,
				Pattern:  "src|lib",
				Original: "*(src|lib)",
			},
		},
		{
			pattern: "@(*.go|*.mod)",
			expected: &ExtglobPattern{
				Type:     ExtglobExact,
				Pattern:  "*.go|*.mod",
				Original: "@(*.go|*.mod)",
			},
		},
		{
			pattern: "*.go",
			wantNil: true,
		},
		{
			pattern: "invalid(pattern",
			wantNil: true,
		},
	}

	for _, c := range cases {
		result := parseExtglob(c.pattern)
		if c.wantNil {
			if result != nil {
				t.Errorf("parseExtglob(%q) should return nil, got %v", c.pattern, result)
			}
			continue
		}
		if result == nil {
			t.Errorf("parseExtglob(%q) returned nil, expected %v", c.pattern, c.expected)
			continue
		}
		if result.Type != c.expected.Type {
			t.Errorf("parseExtglob(%q) Type = %v, want %v", c.pattern, result.Type, c.expected.Type)
		}
		if result.Pattern != c.expected.Pattern {
			t.Errorf("parseExtglob(%q) Pattern = %q, want %q", c.pattern, result.Pattern, c.expected.Pattern)
		}
		if result.Original != c.expected.Original {
			t.Errorf("parseExtglob(%q) Original = %q, want %q", c.pattern, result.Original, c.expected.Original)
		}
	}
}

func TestMatchExtglobNot(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
	}{
		{"*.md", "README.md", false},
		{"*.md", "main.go", true},
		{"**/*.md", "docs/guide.md", false},
		{"*.js|*.ts", "main.go", true},
		{"*.js|*.ts", "app.js", false},
		{"*.js|*.ts", "types.ts", false},
		{"test*|spec*", "main.go", true},
		{"test*|spec*", "test_main.go", false},
		{"test*|spec*", "spec_helper.go", false},
	}

	for _, c := range cases {
		got := matchExtglobNot(c.pattern, c.filename)
		if got != c.want {
			t.Errorf("matchExtglobNot(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchExtglobExact(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
	}{
		{"*.go|*.mod", "main.go", true},
		{"*.go|*.mod", "go.mod", true},
		{"*.go|*.mod", "README.md", false},
		{"test*|spec*", "test_main.go", true},
		{"test*|spec*", "spec_helper.go", true},
		{"test*|spec*", "main.go", false},
	}

	for _, c := range cases {
		got := matchExtglobExact(c.pattern, c.filename)
		if got != c.want {
			t.Errorf("matchExtglobExact(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchGlobWithExtglob(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
	}{
		// Standard glob patterns (should work as before)
		{"*.go", "main.go", true},
		{"*.go", "README.md", false},
		{"docs/*", "docs/readme.md", true},
		{"**/*.go", "src/main.go", true},

		// Extended glob patterns
		{"!(*.md)", "main.go", true},
		{"!(*.md)", "README.md", false},
		{"!(*.js|*.ts)", "main.go", true},
		{"!(*.js|*.ts)", "app.js", false},
		{"@(*.go|*.mod)", "main.go", true},
		{"@(*.go|*.mod)", "go.mod", true},
		{"@(*.go|*.mod)", "README.md", false},

		// Complex patterns with doublestar and extglob
		{"!(docs/**)", "src/main.go", true},
		{"!(docs/**)", "docs/guide.md", false},
		{"!(**/*_test.go)", "main.go", true},
		{"!(**/*_test.go)", "src/main_test.go", false},
	}

	for _, c := range cases {
		got := matchGlob(c.pattern, c.filename)
		if got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestIsExtglob(t *testing.T) {
	cases := []struct {
		pattern string
		want    bool
	}{
		{"!(*.md)", true},
		{"?(*.js)", true},
		{"+(test*)", true},
		{"*(src/*)", true},
		{"@(*.go|*.mod)", true},
		{"*.go", false},
		{"**/*.js", false},
		{"docs/*", false},
		{"invalid(pattern", false},
		{"(*.go)", false}, // Missing extglob prefix
	}

	for _, c := range cases {
		got := isExtglob(c.pattern)
		if got != c.want {
			t.Errorf("isExtglob(%q) = %v, want %v", c.pattern, got, c.want)
		}
	}
}

func TestExtglobType_String(t *testing.T) {
	cases := []struct {
		extType ExtglobType
		want    string
	}{
		{ExtglobNot, "not"},
		{ExtglobZeroOrOne, "zero-or-one"},
		{ExtglobOneOrMore, "one-or-more"},
		{ExtglobZeroOrMore, "zero-or-more"},
		{ExtglobExact, "exact"},
		{ExtglobNone, "none"},
	}

	for _, c := range cases {
		got := c.extType.String()
		if got != c.want {
			t.Errorf("ExtglobType(%v).String() = %q, want %q", c.extType, got, c.want)
		}
	}
}
