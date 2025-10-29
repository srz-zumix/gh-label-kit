package labeler

import (
	"testing"
)

func TestBasicExtglobPatterns(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
		desc     string
	}{
		// Test basic negation
		{
			pattern:  "!(ab)",
			filename: "ab",
			want:     false,
			desc:     "negation should not match exact string",
		},
		{
			pattern:  "!(ab)",
			filename: "cd",
			want:     true,
			desc:     "negation should match different string",
		},
		// Test alternation
		{
			pattern:  "!(ab|def)",
			filename: "ab",
			want:     false,
			desc:     "negation with alternation should not match ab",
		},
		{
			pattern:  "!(ab|def)",
			filename: "def",
			want:     false,
			desc:     "negation with alternation should not match def",
		},
		{
			pattern:  "!(ab|def)",
			filename: "xyz",
			want:     true,
			desc:     "negation with alternation should match other string",
		},
		// Test positive patterns
		{
			pattern:  "+(ab|def)",
			filename: "ab",
			want:     true,
			desc:     "one-or-more should match ab",
		},
		{
			pattern:  "+(ab|def)",
			filename: "def",
			want:     true,
			desc:     "one-or-more should match def",
		},
		{
			pattern:  "+(ab|def)",
			filename: "abdef",
			want:     true,
			desc:     "one-or-more should match abdef",
		},
		{
			pattern:  "+(ab|def)",
			filename: "xyz",
			want:     false,
			desc:     "one-or-more should not match other string",
		},
		// Test wildcard patterns
		{
			pattern:  "@(*.js|*.ts)",
			filename: "app.js",
			want:     true,
			desc:     "exact match should match .js files",
		},
		{
			pattern:  "@(*.js|*.ts)",
			filename: "types.ts",
			want:     true,
			desc:     "exact match should match .ts files",
		},
		{
			pattern:  "@(*.js|*.ts)",
			filename: "main.go",
			want:     false,
			desc:     "exact match should not match .go files",
		},
		// Zero-or-one patterns (simple cases)
		{
			pattern:  "?(test_)*",
			filename: "main.go",
			want:     true,
			desc:     "Should match files with optional test_ prefix",
		},
		{
			pattern:  "?(test_)*",
			filename: "test_helper.go",
			want:     true,
			desc:     "Should match files with test_ prefix",
		},
		{
			pattern:  "?(docs/)*.md",
			filename: "docs/test_helper.go",
			want:     false,
			desc:     "Should not match files with docs/ prefix",
		},
		// Basic extglob with * patterns
		{
			pattern:  "!(test)/*",
			filename: "src/main.go",
			want:     true,
			desc:     "Should match files in non-test directories",
		},
		{
			pattern:  "!(test)/*",
			filename: "test/main.go",
			want:     false,
			desc:     "Should not match files in test directory",
		},
		{
			pattern:  "!(test)/*",
			filename: "src/common/main.go",
			want:     false,
			desc:     "Should not match files in nested directories (standard glob: * excludes /)",
		},
		// Basic extglob with ** patterns
		{
			pattern:  "!(test)/**",
			filename: "src/main.go",
			want:     true,
			desc:     "Should match files in non-test directories",
		},
		{
			pattern:  "!(test)/**",
			filename: "test/helper.go",
			want:     false,
			desc:     "Should not match files in test directory",
		},
		{
			pattern:  "+(src|lib)/**",
			filename: "src/components/Button.js",
			want:     true,
			desc:     "Should match files in src with **",
		},
		{
			pattern:  "+(src|lib)/**",
			filename: "lib/utils/helper.ts",
			want:     true,
			desc:     "Should match files in lib with **",
		},
		{
			pattern:  "+(src|lib)/**",
			filename: "docs/README.md",
			want:     false,
			desc:     "Should not match files outside src/lib with **",
		},
		{
			pattern:  "@(config|docs)/**/*.json",
			filename: "config/app.json",
			want:     true,
			desc:     "Should match JSON files in config directory",
		},
		{
			pattern:  "@(config|docs)/**/*.json",
			filename: "docs/api/schema.json",
			want:     true,
			desc:     "Should match nested JSON files in docs directory",
		},
		{
			pattern:  "@(config|docs)/**/*.json",
			filename: "src/data.json",
			want:     false,
			desc:     "Should not match JSON files outside config/docs",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			// Test with debug output
			t.Logf("Testing pattern: %s with filename: %s (want: %v)", c.pattern, c.filename, c.want)

			// Test the conversion function directly
			regex, ok := convertExtglobToRegex2(c.pattern)
			t.Logf("Converted regex: %s (ok: %v)", regex, ok)

			got := matchGlob(c.pattern, c.filename)
			if got != c.want {
				t.Errorf("%s: matchGlob(%q, %q) = %v, want %v", c.desc, c.pattern, c.filename, got, c.want)
			}
		})
	}
}
