package labeler

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v73/github"
)

func TestCheckMatchConfigs(t *testing.T) {
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr("glob")}}

	testCases := []struct {
		name   string
		config LabelerConfig
		pr     *github.PullRequest
		files  []*github.CommitFile
		want   map[string]bool
	}{
		{
			name: "BranchAndFiles_Any",
			config: LabelerConfig{
				"label1": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}, BaseBranch: []any{"regexp"}, HeadBranch: []any{"regexp"}}}}}},
				"label2": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"hoge"}}}, BaseBranch: []any{"base-branch"}, HeadBranch: []any{"regexp"}}}}}},
				"label3": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"hoge"}}}, BaseBranch: []any{"regexp"}, HeadBranch: []any{"head-branch"}}}}}},
			},
			pr:    pr,
			files: files,
			want: map[string]bool{
				"label1": true,
				"label2": true,
				"label3": true,
			},
		},
		{
			name: "BranchAndFiles_AnyArray",
			config: LabelerConfig{
				"label1": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}}}}, {Any: []LabelerRule{{HeadBranch: []any{"regexp"}}}}, {Any: []LabelerRule{{BaseBranch: []any{"regexp"}}}}}},
				"label2": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}}}}, {Any: []LabelerRule{{HeadBranch: []any{"head-branch"}}}}, {Any: []LabelerRule{{BaseBranch: []any{"base-branch"}}}}}},
			},
			pr:    pr,
			files: files,
			want: map[string]bool{
				"label1": false,
				"label2": true,
			},
		},
		{
			name: "ColorOnly",
			config: LabelerConfig{
				"label1": {Color: "ff0000"},
			},
			pr:    pr,
			files: files,
			want: map[string]bool{
				"label1": false,
			},
		},
		{
			name: "RegexAndDotOption",
			config: LabelerConfig{
				"dotlabel": {Matcher: []LabelerMatch{{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"*.txt"}}}}}}}},
			},
			pr: &github.PullRequest{
				Base:   &github.PullRequestBranch{Ref: github.Ptr("main")},
				Head:   &github.PullRequestBranch{Ref: github.Ptr("feature/abc")},
				Labels: []*github.Label{},
			},
			files: []*github.CommitFile{{Filename: github.Ptr(".foo.txt")}},
			want: map[string]bool{
				"dotlabel": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckMatchConfigs(tc.config, tc.files, tc.pr)
			got := make(map[string]bool)
			for label := range tc.config {
				got[label] = result.IsMatched(label)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("CheckMatchConfigs() = %v, want %v", got, tc.want)
			}
		})
	}
}
