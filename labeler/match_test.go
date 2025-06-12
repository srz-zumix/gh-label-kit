package labeler

import (
	"testing"

	"github.com/google/go-github/v71/github"
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

func TestMatchAnyRegex(t *testing.T) {
	cases := []struct {
		patterns []string
		branch   string
		want     bool
	}{
		{[]string{"^feature"}, "feature/abc", true},
		{[]string{"^hotfix"}, "feature/abc", false},
		{[]string{"^f.*e$"}, "feature", true},
		// 否定先読み
		{[]string{"^(?!main$).*"}, "dev", true},
		{[]string{"^(?!main$).*"}, "main", false},
	}
	for _, c := range cases {
		if got := MatchAnyRegex(c.patterns, c.branch); got != c.want {
			t.Errorf("MatchAnyRegex(%v, %q) = %v, want %v", c.patterns, c.branch, got, c.want)
		}
	}
}

func TestMatchChangedFilesRule(t *testing.T) {
	files := []*github.CommitFile{
		{Filename: strPtr("main.go")},
		{Filename: strPtr("docs/readme.md")},
	}
	cf := ChangedFilesRule{
		AnyGlobToAnyFile: []string{"*.go", "docs/*"},
	}
	if !MatchChangedFilesRule(cf, files) {
		t.Error("AnyGlobToAnyFile should match")
	}
	cf = ChangedFilesRule{
		AnyGlobToAllFiles: []string{"*.go", "docs/*"},
	}
	if MatchChangedFilesRule(cf, files) {
		t.Error("AnyGlobToAllFiles should not match (not all files match any glob)")
	}
	cf = ChangedFilesRule{
		AllGlobsToAnyFile: []string{"*.go", "docs/*"},
	}
	if !MatchChangedFilesRule(cf, files) {
		t.Error("AllGlobsToAnyFile should match (each glob matches at least one file)")
	}
	cf = ChangedFilesRule{
		AllGlobsToAllFiles: []string{"*.go", "docs/*"},
	}
	if MatchChangedFilesRule(cf, files) {
		t.Error("AllGlobsToAllFiles should not match (not all globs match all files)")
	}
}

func strPtr(s string) *string { return &s }
