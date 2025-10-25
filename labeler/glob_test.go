package labeler

import (
	"testing"
)

func TestMatchGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"*.go", "main.go", true},
		{"*.md", "main.go", false},
		{"docs/*", "docs/readme.md", true},
		{"docs/*", "src/readme.md", false},
		{"**", ".github/labeler.yml", true},
		{"*", ".github/labeler.yml", false},
		{"*", "aqua.yml", true},
		{"docs/**", "docs/readme.md", true},
		{"docs/**", "docs/subdir/file.txt", true},
		{"docs/**", "src/readme.md", false},
		{"src/**/test.go", "src/a/test.go", true},
		{"src/**/test.go", "src/test.go", true},
		{"src/**/test.go", "src/a/b/test.go", true},
		{"src/**/test.go", "src/a/b/test.txt", false},
	}
	for _, c := range cases {
		t.Run(c.pattern, func(t *testing.T) {
			if got := matchGlob(c.pattern, c.filename); got != c.want {
				t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
			}
		})
	}
}
