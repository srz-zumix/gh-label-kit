package labeler

import (
	"context"
	"testing"

	"github.com/google/go-github/v79/github"
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
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
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
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
	if result.IsMatched("label1") {
		t.Errorf("label1 should not be matched")
	}
	if !result.IsMatched("label2") {
		t.Errorf("label2 should be matched")
	}
}

func TestCheckMatchConfigs_BranchAndFiles_All(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{All: []LabelerRule{{
					ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}},
					BaseBranch:   []any{"base-branch"},
					HeadBranch:   []any{"base-branch"},
				}}},
			},
		},
		"label2": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{All: []LabelerRule{{
					ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}},
					BaseBranch:   []any{"base-branch"},
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
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
	if result.IsMatched("label1") {
		t.Errorf("label1 should not be matched")
	}
	if !result.IsMatched("label2") {
		t.Errorf("label2 should be matched")
	}
}

func TestCheckMatchConfigs_BranchAndFiles_ColorOnly(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Color: "ff0000",
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr("glob")}}
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
	if result.IsMatched("label1") {
		t.Errorf("label1 should not be matched")
	}
}

func TestCheckMatchConfigs_BranchAndFiles_DescriptionOnly(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Description: "Test label description",
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{{Filename: github.Ptr("glob")}}
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
	if result.IsMatched("label1") {
		t.Errorf("label1 should not be matched")
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
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)

	if !result.IsMatched("dotlabel") {
		t.Errorf("dotlabel should be matched for dotfile")
	}
}
