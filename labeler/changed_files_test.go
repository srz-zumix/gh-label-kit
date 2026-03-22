package labeler

import (
	"testing"
)

func TestMatchChangedFilesRule(t *testing.T) {
	files := []*CommitFile{
		{Filename: Ptr("main.go")},
		{Filename: Ptr("docs/readme.md")},
	}
	cf := ChangedFilesRule{
		AnyGlobToAnyFile: []string{"*.go", "docs/*"},
	}
	if !matchChangedFilesRuleAny(cf, files) {
		t.Error("AnyGlobToAnyFile should match")
	}
	cf = ChangedFilesRule{
		AnyGlobToAllFiles: []string{"*.go", "docs/*"},
	}
	if matchChangedFilesRuleAny(cf, files) {
		t.Error("AnyGlobToAllFiles should not match (not all files match any glob)")
	}
	cf = ChangedFilesRule{
		AllGlobsToAnyFile: []string{"*.go", "docs/*"},
	}
	if matchChangedFilesRuleAny(cf, files) {
		t.Error("AllGlobsToAnyFile should not match (no single file matches all globs)")
	}
	// positive case: main.go matches both *.go and main*
	cf = ChangedFilesRule{
		AllGlobsToAnyFile: []string{"*.go", "main*"},
	}
	if !matchChangedFilesRuleAny(cf, files) {
		t.Error("AllGlobsToAnyFile should match (main.go matches all globs)")
	}
	cf = ChangedFilesRule{
		AllGlobsToAllFiles: []string{"*.go", "docs/*"},
	}
	if matchChangedFilesRuleAny(cf, files) {
		t.Error("AllGlobsToAllFiles should not match (not all globs match all files)")
	}
}

func TestCheckAllChangedFiles(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	// all configs matched
	configs := []ChangedFilesRule{
		{AnyGlobToAnyFile: []string{"foo.txt"}},
		{AnyGlobToAllFiles: []string{"*.txt"}},
		{AllGlobsToAllFiles: []string{"**"}},
	}
	allMatch := matchChangedFilesAll(configs, changedFiles)
	if !allMatch {
		t.Error("all configs should match")
	}
	// some configs not matched
	configs = []ChangedFilesRule{
		{AnyGlobToAnyFile: []string{"foo.txt"}},
		{AnyGlobToAllFiles: []string{"*.md"}},
		{AllGlobsToAllFiles: []string{"**"}},
	}
	allMatch = matchChangedFilesAll(configs, changedFiles)
	if allMatch {
		t.Error("not all configs should match")
	}
}

func TestCheckAnyChangedFiles(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	// any config matched
	configs := []ChangedFilesRule{
		{AnyGlobToAnyFile: []string{"*.md"}},
		{AnyGlobToAllFiles: []string{"*.txt"}},
	}
	anyMatch := matchChangedFilesAny(configs, changedFiles)
	if !anyMatch {
		t.Error("at least one config should match")
	}
	// none matched
	configs = []ChangedFilesRule{
		{AnyGlobToAnyFile: []string{"*.md"}},
		{AnyGlobToAllFiles: []string{"!*.txt"}},
	}
	anyMatch = matchChangedFilesAny(configs, changedFiles)
	if anyMatch {
		t.Error("no config should match")
	}
}

func TestMatchChangedFilesRule_AllGlobsMatchAnyFile(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	cf := ChangedFilesRule{AllGlobsToAnyFile: []string{"**/bar.txt", "bar.txt"}}
	if !matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("all globs should match any file")
	}
	cf = ChangedFilesRule{AllGlobsToAnyFile: []string{"*.txt", "*.md"}}
	if matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("not all globs should match any file")
	}
}

func TestMatchChangedFilesRule_AnyGlobMatchesAllFiles(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	cf := ChangedFilesRule{AnyGlobToAllFiles: []string{"*.md", "*.txt"}}
	if !matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("any glob should match all files")
	}
	cf = ChangedFilesRule{AnyGlobToAllFiles: []string{"*.md", "bar.txt", "foo.txt"}}
	if matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("no glob should match all files")
	}
}

func TestMatchChangedFilesRule_AllGlobsMatchAllFiles(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	cf := ChangedFilesRule{AllGlobsToAllFiles: []string{"*.txt", "**"}}
	if !matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("all globs should match all files")
	}
	cf = ChangedFilesRule{AllGlobsToAllFiles: []string{"**", "foo.txt"}}
	if matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("not all globs should match all files")
	}
}

func TestMatchChangedFilesRule_AllFilesToAnyGlob(t *testing.T) {
	changedFiles := []*CommitFile{
		{Filename: Ptr("foo.txt")},
		{Filename: Ptr("bar.txt")},
	}
	// all files match at least one glob pattern
	cf := ChangedFilesRule{AllFilesToAnyGlob: []string{"*.txt", "*.md"}}
	if !matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("all files should match at least one glob")
	}
	// not all files match any glob pattern
	cf = ChangedFilesRule{AllFilesToAnyGlob: []string{"foo.txt", "*.md"}}
	if matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("bar.txt does not match any glob")
	}
	// all files match when using wildcard
	cf = ChangedFilesRule{AllFilesToAnyGlob: []string{"**"}}
	if !matchChangedFilesRuleAll(cf, changedFiles) {
		t.Error("all files should match wildcard")
	}
}

// TestMatchFunctions_Comprehensive tests all five match functions with shared test cases
// to ensure each function has distinct and correct semantics.
//
// Semantics:
//
//	matchAnyGlobToAnyFile:  ∃glob ∃file: match(glob, file)
//	matchAnyGlobToAllFiles: ∃glob ∀file: match(glob, file)
//	matchAllGlobsToAnyFile: ∃file ∀glob: match(glob, file)
//	matchAllGlobsToAllFiles: ∀glob ∀file: match(glob, file)
//	matchAllFilesToAnyGlob: ∀file ∃glob: match(glob, file)
func TestMatchFunctions_Comprehensive(t *testing.T) {
	tests := []struct {
		name                   string
		patterns               []string
		files                  []string
		wantAnyGlobToAnyFile   bool
		wantAnyGlobToAllFiles  bool
		wantAllGlobsToAnyFile  bool
		wantAllGlobsToAllFiles bool
		wantAllFilesToAnyGlob  bool
	}{
		{
			// Key case: each glob matches a different file, but no single file matches all globs
			name:                   "different globs match different files",
			patterns:               []string{"*.go", "*.txt"},
			files:                  []string{"main.go", "foo.txt"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  false,
			wantAllGlobsToAnyFile:  false, // no single file matches both *.go AND *.txt
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  true, // each file matches at least one glob
		},
		{
			// A single file matches all globs
			name:                   "single file matches all globs",
			patterns:               []string{"*.go", "main.*"},
			files:                  []string{"main.go"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  true,
			wantAllGlobsToAnyFile:  true, // main.go matches both *.go and main.*
			wantAllGlobsToAllFiles: true,
			wantAllFilesToAnyGlob:  true,
		},
		{
			// Multiple files, one file matches all globs
			name:                   "one of multiple files matches all globs",
			patterns:               []string{"*.go", "main.*"},
			files:                  []string{"main.go", "test.go"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  true,  // *.go matches both main.go and test.go
			wantAllGlobsToAnyFile:  true,  // main.go matches both *.go and main.*
			wantAllGlobsToAllFiles: false, // main.* doesn't match test.go
			wantAllFilesToAnyGlob:  true,  // main.go→*.go, test.go→*.go
		},
		{
			// Wildcard matches everything
			name:                   "wildcard matches all",
			patterns:               []string{"**"},
			files:                  []string{"main.go", "foo.txt"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  true,
			wantAllGlobsToAnyFile:  true,
			wantAllGlobsToAllFiles: true,
			wantAllFilesToAnyGlob:  true,
		},
		{
			// No glob matches any file
			name:                   "no match at all",
			patterns:               []string{"*.md"},
			files:                  []string{"main.go", "foo.txt"},
			wantAnyGlobToAnyFile:   false,
			wantAnyGlobToAllFiles:  false,
			wantAllGlobsToAnyFile:  false,
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  false,
		},
		{
			// One glob matches all files, another matches none
			name:                   "one glob matches all, another matches none",
			patterns:               []string{"*.go", "*.md"},
			files:                  []string{"main.go", "test.go"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  true,  // *.go matches all files
			wantAllGlobsToAnyFile:  false, // no file matches both *.go AND *.md
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  true, // main.go→*.go ✓, test.go→*.go ✓ (each file needs at least one)
		},
		{
			// Three files, three globs, each glob matches a different file
			name:                   "three globs three files, no single file matches all",
			patterns:               []string{"*.go", "*.md", "*.txt"},
			files:                  []string{"main.go", "README.md", "notes.txt"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  false,
			wantAllGlobsToAnyFile:  false, // no single file matches all three globs
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  true, // each file matches its corresponding glob
		},
		{
			// Empty files
			name:                   "empty files",
			patterns:               []string{"*.go"},
			files:                  []string{},
			wantAnyGlobToAnyFile:   false,
			wantAnyGlobToAllFiles:  false,
			wantAllGlobsToAnyFile:  false,
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  false,
		},
		{
			// Single glob, single file, match
			name:                   "single glob single file match",
			patterns:               []string{"*.go"},
			files:                  []string{"main.go"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  true,
			wantAllGlobsToAnyFile:  true,
			wantAllGlobsToAllFiles: true,
			wantAllFilesToAnyGlob:  true,
		},
		{
			// Single glob, single file, no match
			name:                   "single glob single file no match",
			patterns:               []string{"*.md"},
			files:                  []string{"main.go"},
			wantAnyGlobToAnyFile:   false,
			wantAnyGlobToAllFiles:  false,
			wantAllGlobsToAnyFile:  false,
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  false,
		},
		{
			// Partial overlap: some files match all globs, some don't
			name:                   "partial overlap with directory patterns",
			patterns:               []string{"src/**", "**/*.go"},
			files:                  []string{"src/main.go", "README.md"},
			wantAnyGlobToAnyFile:   true,
			wantAnyGlobToAllFiles:  false, // src/** doesn't match README.md
			wantAllGlobsToAnyFile:  true,  // src/main.go matches both src/** and **/*.go
			wantAllGlobsToAllFiles: false,
			wantAllFilesToAnyGlob:  false, // README.md doesn't match either
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changedFiles := make([]*CommitFile, len(tt.files))
			for i, f := range tt.files {
				changedFiles[i] = &CommitFile{Filename: Ptr(f)}
			}

			if got := matchAnyGlobToAnyFile(tt.patterns, changedFiles); got != tt.wantAnyGlobToAnyFile {
				t.Errorf("matchAnyGlobToAnyFile() = %v, want %v", got, tt.wantAnyGlobToAnyFile)
			}
			if got := matchAnyGlobToAllFiles(tt.patterns, changedFiles); got != tt.wantAnyGlobToAllFiles {
				t.Errorf("matchAnyGlobToAllFiles() = %v, want %v", got, tt.wantAnyGlobToAllFiles)
			}
			if got := matchAllGlobsToAnyFile(tt.patterns, changedFiles); got != tt.wantAllGlobsToAnyFile {
				t.Errorf("matchAllGlobsToAnyFile() = %v, want %v", got, tt.wantAllGlobsToAnyFile)
			}
			if got := matchAllGlobsToAllFiles(tt.patterns, changedFiles); got != tt.wantAllGlobsToAllFiles {
				t.Errorf("matchAllGlobsToAllFiles() = %v, want %v", got, tt.wantAllGlobsToAllFiles)
			}
			if got := matchAllFilesToAnyGlob(tt.patterns, changedFiles); got != tt.wantAllFilesToAnyGlob {
				t.Errorf("matchAllFilesToAnyGlob() = %v, want %v", got, tt.wantAllFilesToAnyGlob)
			}
		})
	}
}
