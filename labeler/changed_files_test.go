package labeler

import (
	"testing"

	"github.com/google/go-github/v79/github"
)

func TestMatchChangedFilesRule(t *testing.T) {
	files := []*github.CommitFile{
		{Filename: github.Ptr("main.go")},
		{Filename: github.Ptr("docs/readme.md")},
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
	if !matchChangedFilesRuleAny(cf, files) {
		t.Error("AllGlobsToAnyFile should match (each glob matches at least one file)")
	}
	cf = ChangedFilesRule{
		AllGlobsToAllFiles: []string{"*.go", "docs/*"},
	}
	if matchChangedFilesRuleAny(cf, files) {
		t.Error("AllGlobsToAllFiles should not match (not all globs match all files)")
	}
}

func TestCheckAllChangedFiles(t *testing.T) {
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
	changedFiles := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
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
