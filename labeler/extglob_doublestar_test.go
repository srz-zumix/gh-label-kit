package labeler

import (
	"testing"
)

// TestExtglobWithDoubleStarPatterns tests extglob patterns followed by ** (doublestar) patterns
func TestExtglobWithDoubleStarPatterns(t *testing.T) {
	cases := []struct {
		pattern  string
		filename string
		want     bool
		desc     string
	}{
		// Multiple alternatives with **
		{
			pattern:  "!(test|spec|__tests__)/**",
			filename: "src/components/Button.js",
			want:     true,
			desc:     "Should match files outside test directories",
		},
		{
			pattern:  "!(test|spec|__tests__)/**",
			filename: "lib/utils/helper.ts",
			want:     true,
			desc:     "Should match nested files outside test directories",
		},
		{
			pattern:  "!(test|spec|__tests__)/**",
			filename: "test/unit.js",
			want:     false,
			desc:     "Should not match files in test directory",
		},
		{
			pattern:  "!(test|spec|__tests__)/**",
			filename: "spec/helper.rb",
			want:     false,
			desc:     "Should not match files in spec directory",
		},
		{
			pattern:  "!(test|spec|__tests__)/**",
			filename: "__tests__/component.test.js",
			want:     false,
			desc:     "Should not match files in __tests__ directory",
		},

		// ** followed by extglob
		{
			pattern:  "**/!(*.test.js)",
			filename: "src/main.js",
			want:     true,
			desc:     "Should match non-test JS files with ** prefix",
		},
		{
			pattern:  "**/!(*.test.js)",
			filename: "lib/utils/helper.js",
			want:     true,
			desc:     "Should match nested non-test JS files",
		},
		{
			pattern:  "**/!(*.test.js)",
			filename: "src/component.test.js",
			want:     false,
			desc:     "Should not match test JS files",
		},
		{
			pattern:  "**/!(*.test.js)",
			filename: "lib/nested/unit.test.js",
			want:     false,
			desc:     "Should not match nested test JS files",
		},

		// Complex combinations
		{
			pattern:  "!(node_modules)/**/*.@(js|ts)",
			filename: "src/app.js",
			want:     true,
			desc:     "Should match JS files outside node_modules",
		},
		{
			pattern:  "!(node_modules)/**/*.@(js|ts)",
			filename: "node_modules/react/index.js",
			want:     false,
			desc:     "Should not match JS files in node_modules",
		},

		// Multiple ** with extglob
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "cmd/main/app.go",
			want:     true,
			desc:     "Should match Go files in non-test nested directories",
		},
		{
			pattern:  "**/!(test)/**/*.go",
			filename: "cmd/test/main.go",
			want:     false,
			desc:     "Should not match Go files in test directories",
		},

		// Edge cases
		{
			pattern:  "!(*)/**",
			filename: "any/file.txt",
			want:     false,
			desc:     "!(*)/** should not match anything (negates everything)",
		},
		{
			pattern:  "!()/**",
			filename: "src/main.go",
			want:     true,
			desc:     "!()/** should match everything (negates empty pattern)",
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
