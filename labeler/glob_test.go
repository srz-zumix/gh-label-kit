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

func TestMatchGlob_NegationGlob(t *testing.T) {
	cases := []struct {
		pattern, filename string
		want              bool
	}{
		{"!*.go", "main.go", false},
		{"!*.go", "README.md", true},
		{"!docs/*", "docs/readme.md", false},
		{"!docs/*", "src/readme.md", true},
		{"!**/*.test.js", "src/test.js", true},
		{"!**/*.test.js", "src/utils.test.js", false},
		{"!*.md", "file.txt", true},
		{"!*.md", "file.md", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", c.pattern, c.filename, got, c.want)
		}
	}
}

func TestIsHiddenFile(t *testing.T) {
	cases := []struct {
		filename string
		want     bool
	}{
		// Hidden files (starting with .)
		{".hidden", true},
		{".gitignore", true},
		{".github/workflows/test.yml", true},
		{"dir/.hidden", true},
		{"dir/.config/file.txt", true},
		{"a/b/.c/d.txt", true},

		// Visible files
		{"visible.txt", false},
		{"file.go", false},
		{"src/main.go", false},
		{"docs/readme.md", false},
		{"a/b/c/d.txt", false},

		// Edge cases
		{"", false},
		{".", true},
		{"..", true},
		{"...", true},
		{"...visible", true},
		{"visible.", false},
		{"dir/file.", false},
	}

	for _, c := range cases {
		if got := isHiddenFile(c.filename); got != c.want {
			t.Errorf("isHiddenFile(%q) = %v, want %v", c.filename, got, c.want)
		}
	}
}

func TestMatchGlob_NoHidden(t *testing.T) {
	// Save original state
	original := noHiddenEnabled
	defer func() {
		noHiddenEnabled = original
	}()

	cases := []struct {
		pattern, filename string
		noHidden          bool
		want              bool
	}{
		// Without no-hidden option
		{"**/*.go", "main.go", false, true},
		{"**/*.go", ".hidden.go", false, true},
		{"**/*.yml", ".github/workflows/test.yml", false, true},
		{"**/*", "dir/.config/file.txt", false, true},

		// With no-hidden option
		{"**/*.go", "main.go", true, true},
		{"**/*.go", ".hidden.go", true, false},
		{"**/*.yml", ".github/workflows/test.yml", true, false},
		{"**/*", "dir/.config/file.txt", true, false},
		{"**/*", "visible/file.txt", true, true},
		{"src/**/*.go", "src/main.go", true, true},
		{"src/**/*.go", "src/.hidden.go", true, false},
		{"src/**/*.go", "src/pkg/file.go", true, true},
		{"src/**/*.go", "src/.internal/file.go", true, false},
	}

	for _, c := range cases {
		SetNoHidden(c.noHidden)
		if got := matchGlob(c.pattern, c.filename); got != c.want {
			t.Errorf("matchGlob(%q, %q) with noHidden=%v = %v, want %v",
				c.pattern, c.filename, c.noHidden, got, c.want)
		}
	}
}
