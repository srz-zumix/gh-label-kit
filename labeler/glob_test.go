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
		{"ab*.jpg", "abc.jpg", true},
		{"ab*.jpg", "abc/test.jpg", false},
		{"**", ".github/labeler.yml", true},
		{"*", ".github/labeler.yml", false},
		{"*", "aqua.yml", true},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchGlob_DirectoryGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"docs/**", "docs/readme.md", true},
		{"docs/**", "docs/subdir/file.txt", true},
		{"docs/**", "src/readme.md", false},
		{"src/**/test.go", "src/a/test.go", true},
		{"src/**/test.go", "src/test.go", true},
		{"src/**/test.go", "src/a/b/test.go", true},
		{"src/**/test.go", "src/a/b/test.txt", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}
