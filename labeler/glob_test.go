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

func TestMatchGlob_SingleCharacterGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"file?.txt", "file1.txt", true},
		{"file?.txt", "fileA.txt", true},
		{"file?.txt", "file12.txt", false},
		{"data-??.csv", "data-01.csv", true},
		{"data-??.csv", "data-1.csv", false},
		{"data-??.csv", "data-001.csv", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchGlob_CharacterClassGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"file[0-9].txt", "file1.txt", true},
		{"file[0-9].txt", "fileA.txt", false},
		{"file[0-9].txt", "file12.txt", false},
		{"file[^0-9].txt", "fileA.txt", true},
		{"file[!0-9].txt", "fileA.txt", true},
		{"data-[0-9][0-9].csv", "data-01.csv", true},
		{"data-[0-9][0-9].csv", "data-1.csv", false},
		{"data-[0-9][0-9].csv", "data-001.csv", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchGlob_AlternativesGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"docs/some{thing{new,},}.md", "docs/readme.md", false},
		{"docs/some{thing{new,},}.md", "docs/somenew.md", false},
		{"docs/some{thing{new,},}.md", "docs/something.md", true},
		{"docs/some{thing{new,},}.md", "docs/somethingnew.md", true},
		{"docs/some{thing{new,},}.md", "docs/some.md", true},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestMatchGlob_CaseSensitiveGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"README.md", "README.md", true},
		{"README.md", "readme.md", false},
		{"CMake", "cmake", false},
		{"CMake", "CMake", true},
		{"[cC][mM]ake", "cmake", true},
		{"[cC][mM]ake", "CMake", true},
		{"[cC][mM]ake", "Cmake", true},
		{"{CM,cm}ake", "cmake", true},
		{"{CM,cm}ake", "CMake", true},
		{"{CM,cm}ake", "Cmake", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}
