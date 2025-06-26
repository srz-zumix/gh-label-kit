package labeler

import (
	"testing"

	"github.com/google/go-github/v71/github"
)

func TestCheckMatchConfigs_BranchAndFiles_Any(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}},
					BaseBranch:   []any{"regexp"},
					HeadBranch:   []any{"regexp"},
				}}},
			},
		},
		"label2": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"hoge"}}},
					BaseBranch:   []any{"base-branch"},
					HeadBranch:   []any{"regexp"},
				}}},
			},
		},
		"label3": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"hoge"}}},
					BaseBranch:   []any{"regexp"},
					HeadBranch:   []any{"head-branch"},
				}}},
			},
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr("glob")}}
	result := CheckMatchConfigs(cfg, files, pr)
	if !result.IsMatched("label1") {
		t.Errorf("label1 should be matched")
	}
	if !result.IsMatched("label2") {
		t.Errorf("label2 should be matched")
	}
	if !result.IsMatched("label3") {
		t.Errorf("label3 should be matched")
	}
}

func TestCheckMatchConfigs_BranchAndFiles_AnyArray(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}}}},
				{Any: []LabelerRule{{HeadBranch: []any{"regexp"}}}},
				{Any: []LabelerRule{{BaseBranch: []any{"regexp"}}}},
			},
		},
		"label2": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}}}},
				{Any: []LabelerRule{{HeadBranch: []any{"head-branch"}}}},
				{Any: []LabelerRule{{BaseBranch: []any{"base-branch"}}}},
			},
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr("glob")}}
	result := CheckMatchConfigs(cfg, files, pr)
	if result.IsMatched("label1") {
		t.Errorf("label1 should not be matched")
	}
	if !result.IsMatched("label2") {
		t.Errorf("label2 should be matched")
	}
}

func TestCheckMatchConfigs_RegexAndDotOption(t *testing.T) {
	cfg := LabelerConfig{
		"dotlabel": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"*.txt"}}}}}},
			},
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("main")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("feature/abc")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr(".foo.txt")}}
	// matchGlob uses doublestar, which matches dotfiles by default
	result := CheckMatchConfigs(cfg, files, pr)

	if !result.IsMatched("dotlabel") {
		t.Errorf("dotlabel should be matched for dotfile")
	}
}
