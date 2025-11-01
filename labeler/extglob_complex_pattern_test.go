package labeler

import (
	"testing"
)

func TestComplexExtglobPatterns(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
		desc     string
	}{
		// Nested extglob patterns (working cases)
		{
			pattern:  "!(+(ab|def)*)",
			filename: "ab",
			want:     false,
			desc:     "Should not match 'ab' (matches +(ab|def)*)",
		},
		{
			pattern:  "!(+(ab|def)*)",
			filename: "abab",
			want:     false,
			desc:     "Should not match 'abab' (matches +(ab|def)*)",
		},
		{
			pattern:  "!(+(ab|def)*)",
			filename: "def",
			want:     false,
			desc:     "Should not match 'def' (matches +(ab|def)*)",
		},
		{
			pattern:  "!(+(ab|def)*)",
			filename: "defdef",
			want:     false,
			desc:     "Should not match 'defdef' (matches +(ab|def)*)",
		},
		{
			pattern:  "!(+(ab|def)*)",
			filename: "xyz",
			want:     true,
			desc:     "Should match 'xyz' (does not match +(ab|def)*)",
		},
		{
			pattern:  "!(+(ab|def)*)",
			filename: "",
			want:     true,
			desc:     "Should match empty string (does not match +(ab|def)*)",
		},
		// Nested with file extensions (working cases)
		{
			pattern:  "@(*.@(js|ts))",
			filename: "app.js",
			want:     true,
			desc:     "Should match .js files in nested pattern",
		},
		{
			pattern:  "@(*.@(js|ts))",
			filename: "types.ts",
			want:     true,
			desc:     "Should match .ts files in nested pattern",
		},
		{
			pattern:  "@(*.@(js|ts))",
			filename: "main.go",
			want:     false,
			desc:     "Should not match .go files in nested pattern",
		},
		{
			pattern:  "@(*.@(js|ts))",
			filename: "app.jsx",
			want:     false,
			desc:     "Should not match .jsx files in nested pattern",
		},
		// Simple negation patterns (current limitations)
		{
			pattern:  "!(*.md)",
			filename: "main.go",
			want:     true,
			desc:     "Should match non-markdown files",
		},
		{
			pattern:  "!(*.md)",
			filename: "README.md",
			want:     false,
			desc:     "Should not match markdown files",
		},
		{
			pattern:  "!(test*)",
			filename: "src/main.go",
			want:     true,
			desc:     "Should match files not starting with test",
		},
		{
			pattern:  "!(test*)",
			filename: "test_helper.go",
			want:     false,
			desc:     "Should not match files starting with test",
		},
		// Combination patterns
		{
			pattern:  "+(src|lib)/@(*.@(js|ts))",
			filename: "src/app.js",
			want:     true,
			desc:     "Should match JS files in src directory",
		},
		{
			pattern:  "+(src|lib)/@(*.@(js|ts))",
			filename: "lib/utils.ts",
			want:     true,
			desc:     "Should match TS files in lib directory",
		},
		{
			pattern:  "+(src|lib)/@(*.@(js|ts))",
			filename: "docs/README.md",
			want:     false,
			desc:     "Should not match files in docs directory",
		},
		{
			pattern:  "+(src|lib)/@(*.@(js|ts))",
			filename: "src/app.go",
			want:     false,
			desc:     "Should not match Go files even in src directory",
		},
		// Extglob with ** patterns (directory traversal)
		{
			pattern:  "!(test|spec)/**",
			filename: "src/main.go",
			want:     true,
			desc:     "Should match files in non-test directories with **",
		},
		{
			pattern:  "!(test|spec)/**",
			filename: "lib/utils/helper.go",
			want:     true,
			desc:     "Should match nested files in non-test directories",
		},
		{
			pattern:  "!(test|spec)/**",
			filename: "test/main_test.go",
			want:     false,
			desc:     "Should not match files in test directory",
		},
		{
			pattern:  "!(test|spec)/**",
			filename: "spec/helper.go",
			want:     false,
			desc:     "Should not match files in spec directory",
		},
		{
			pattern:  "!(test|spec)/**",
			filename: "test/nested/deep_test.go",
			want:     false,
			desc:     "Should not match deeply nested files in test directory",
		},
		{
			pattern:  "!(test|spec)/**",
			filename: "spec/unit/helper_spec.go",
			want:     false,
			desc:     "Should not match deeply nested files in spec directory",
		},
		// More complex ** patterns
		{
			pattern:  "!(node_modules)/**/*.js",
			filename: "src/components/Button.js",
			want:     true,
			desc:     "Should match JS files outside node_modules",
		},
		{
			pattern:  "!(node_modules)/**/*.js",
			filename: "lib/utils/index.js",
			want:     true,
			desc:     "Should match nested JS files outside node_modules",
		},
		{
			pattern:  "!(node_modules)/**/*.js",
			filename: "node_modules/react/index.js",
			want:     false,
			desc:     "Should not match JS files in node_modules",
		},
		{
			pattern:  "!(node_modules)/**/*.js",
			filename: "node_modules/lodash/lib/core.js",
			want:     false,
			desc:     "Should not match deeply nested JS files in node_modules",
		},
		// Extglob at the end with ** in middle
		{
			pattern:  "src/**/@(*.test.@(js|ts))",
			filename: "src/components/Button.test.js",
			want:     true,
			desc:     "Should match test JS files in src with ** traversal",
		},
		{
			pattern:  "src/**/@(*.test.@(js|ts))",
			filename: "src/utils/nested/helper.test.ts",
			want:     true,
			desc:     "Should match deeply nested test TS files in src",
		},
		{
			pattern:  "src/**/@(*.test.@(js|ts))",
			filename: "lib/utils/helper.test.js",
			want:     false,
			desc:     "Should not match test files outside src directory",
		},
		{
			pattern:  "src/**/@(*.test.@(js|ts))",
			filename: "src/components/Button.js",
			want:     false,
			desc:     "Should not match non-test JS files in src",
		},
		// Multiple ** with extglob
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "cmd/main/main.go",
			want:     true,
			desc:     "Should match Go files in non-test directories with multiple **",
		},
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "pkg/utils/helper.go",
			want:     true,
			desc:     "Should match Go files in nested non-test directories",
		},
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "cmd/test/test_main.go",
			want:     false,
			desc:     "Should not match Go files in test directories",
		},
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "internal/test/mock.go",
			want:     false,
			desc:     "Should not match Go files in nested test directories",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			got := matchGlob(c.pattern, c.filename)
			if got != c.want {
				t.Errorf("%s: matchGlob(%q, %q) = %v, want %v", c.desc, c.pattern, c.filename, got, c.want)
			}
		})
	}
}
