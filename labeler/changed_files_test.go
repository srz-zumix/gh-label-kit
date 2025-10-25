package labeler

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-github/v73/github"
)

func TestChangedFiles(t *testing.T) {
	files := []*github.CommitFile{
		{Filename: github.Ptr("main.go")},
		{Filename: github.Ptr("docs/readme.md")},
	}
	filesTxt := []*github.CommitFile{
		{Filename: github.Ptr("foo.txt")},
		{Filename: github.Ptr("bar.txt")},
	}

	testCases := []struct {
		name         string
		rules        []ChangedFilesRule
		changedFiles []*github.CommitFile
		matcher      func([]ChangedFilesRule, []*github.CommitFile) bool
		want         bool
	}{
		{
			name: "AnyGlobToAnyFile should match",
			rules: []ChangedFilesRule{
				{AnyGlobToAnyFile: []string{"*.go", "docs/*"}},
			},
			changedFiles: files,
			matcher:      matchChangedFilesAny,
			want:         true,
		},
		{
			name: "AnyGlobToAllFiles should not match (not all files match any glob)",
			rules: []ChangedFilesRule{
				{AnyGlobToAllFiles: []string{"*.go", "docs/*"}},
			},
			changedFiles: files,
			matcher:      matchChangedFilesAny,
			want:         false,
		},
		{
			name: "AllGlobsToAnyFile should match (each glob matches at least one file)",
			rules: []ChangedFilesRule{
				{AllGlobsToAnyFile: []string{"*.go", "docs/*"}},
			},
			changedFiles: files,
			matcher:      matchChangedFilesAny,
			want:         true,
		},
		{
			name: "AllGlobsToAllFiles should not match (not all globs match all files)",
			rules: []ChangedFilesRule{
				{AllGlobsToAllFiles: []string{"*.go", "docs/*"}},
			},
			changedFiles: files,
			matcher:      matchChangedFilesAny,
			want:         false,
		},
		{
			name: "all configs should match",
			rules: []ChangedFilesRule{
				{AnyGlobToAnyFile: []string{"foo.txt"}},
				{AnyGlobToAllFiles: []string{"*.txt"}},
				{AllGlobsToAllFiles: []string{"**"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         true,
		},
		{
			name: "not all configs should match",
			rules: []ChangedFilesRule{
				{AnyGlobToAnyFile: []string{"foo.txt"}},
				{AnyGlobToAllFiles: []string{"*.md"}},
				{AllGlobsToAllFiles: []string{"**"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         false,
		},
		{
			name: "at least one config should match",
			rules: []ChangedFilesRule{
				{AnyGlobToAnyFile: []string{"*.md"}},
				{AnyGlobToAllFiles: []string{"*.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAny,
			want:         true,
		},
		{
			name: "no config should match",
			rules: []ChangedFilesRule{
				{AnyGlobToAnyFile: []string{"*.md"}},
				{AnyGlobToAllFiles: []string{"!*.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAny,
			want:         false,
		},
		{
			name: "all globs should match any file",
			rules: []ChangedFilesRule{
				{AllGlobsToAnyFile: []string{"**/bar.txt", "bar.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         true,
		},
		{
			name: "not all globs should match any file",
			rules: []ChangedFilesRule{
				{AllGlobsToAnyFile: []string{"*.txt", "*.md"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         false,
		},
		{
			name: "any glob should match all files",
			rules: []ChangedFilesRule{
				{AnyGlobToAllFiles: []string{"*.md", "*.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         true,
		},
		{
			name: "no glob should match all files",
			rules: []ChangedFilesRule{
				{AnyGlobToAllFiles: []string{"*.md", "bar.txt", "foo.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         false,
		},
		{
			name: "all globs should match all files",
			rules: []ChangedFilesRule{
				{AllGlobsToAllFiles: []string{"*.txt", "**"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         true,
		},
		{
			name: "not all globs should match all files",
			rules: []ChangedFilesRule{
				{AllGlobsToAllFiles: []string{"**", "foo.txt"}},
			},
			changedFiles: filesTxt,
			matcher:      matchChangedFilesAll,
			want:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matcherName := runtime.FuncForPC(reflect.ValueOf(tc.matcher).Pointer()).Name()
			if got := tc.matcher(tc.rules, tc.changedFiles); got != tc.want {
				t.Errorf("%s() = %v, want %v", matcherName, got, tc.want)
			}
		})
	}
}
