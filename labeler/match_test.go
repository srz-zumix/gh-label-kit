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
		"label3": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}}}},
				{Any: []LabelerRule{{HeadBranch: []any{"head-branch"}}}},
				{Any: []LabelerRule{{BaseBranch: []any{"regexp"}}}},
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
	if result.IsMatched("label3") {
		t.Errorf("label3 should not be matched")
	}
}

func TestCheckMatchConfigs_AllFilesToAnyGlob(t *testing.T) {
	cfg := LabelerConfig{
		"label1": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{
					{ChangedFiles: []ChangedFilesRule{{AnyGlobToAnyFile: []string{"glob"}}}},
					{ChangedFiles: []ChangedFilesRule{{AllFilesToAnyGlob: []string{
						".github/**",
						"*.yml",
					}}}},
				}},
			},
		},
	}
	pr := &github.PullRequest{
		Base:   &github.PullRequestBranch{Ref: github.Ptr("base-branch")},
		Head:   &github.PullRequestBranch{Ref: github.Ptr("head-branch")},
		Labels: []*github.Label{},
	}
	files := []*github.CommitFile{
		{Filename: github.Ptr(".github/workflows/test.yml")},
		{Filename: github.Ptr("zizmor.yml")},
	}
	matcher := NewMatcher(context.TODO(), nil)
	result := matcher.CheckMatchConfigs(cfg, files, pr)
	if !result.IsMatched("label1") {
		t.Errorf("label1 should be matched")
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

// TestCheckMatchConfigs_MultipleBranchPatternsInAny tests that multiple branch patterns
// within a single LabelerRule are evaluated as OR (any match succeeds)
// IMPORTANT: When both HeadBranch and BaseBranch are in the same rule under "any",
// they are evaluated as OR - meaning if EITHER matches, the rule succeeds.
func TestCheckMatchConfigs_MultipleBranchPatternsInAny(t *testing.T) {
	// Case 1: Multiple head-branch patterns in a single rule (should be OR)
	cfg := LabelerConfig{
		"multi-head-or": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					HeadBranch: []any{"feature/.*", "bugfix/.*", "hotfix/.*"},
				}}},
			},
		},
		"multi-base-or": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					BaseBranch: []any{"main", "develop", "release/.*"},
				}}},
			},
		},
		// When HeadBranch and BaseBranch are in the same rule under "any",
		// they are OR'd together - if either matches, the rule succeeds
		"multi-both-or": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{
					HeadBranch: []any{"feature/.*", "bugfix/.*"},
					BaseBranch: []any{"main", "develop"},
				}}},
			},
		},
	}

	tests := []struct {
		name          string
		headRef       string
		baseRef       string
		expectMatched []string
	}{
		{
			name:          "head matches first pattern",
			headRef:       "feature/new-feature",
			baseRef:       "main",
			expectMatched: []string{"multi-head-or", "multi-base-or", "multi-both-or"},
		},
		{
			name:          "head matches second pattern",
			headRef:       "bugfix/fix-123",
			baseRef:       "main",
			expectMatched: []string{"multi-head-or", "multi-base-or", "multi-both-or"},
		},
		{
			name:          "head matches third pattern (hotfix), base matches main",
			headRef:       "hotfix/urgent",
			baseRef:       "main",
			// multi-both-or matches because base=main matches BaseBranch patterns
			expectMatched: []string{"multi-head-or", "multi-base-or", "multi-both-or"},
		},
		{
			name:          "base matches develop",
			headRef:       "feature/abc",
			baseRef:       "develop",
			expectMatched: []string{"multi-head-or", "multi-base-or", "multi-both-or"},
		},
		{
			name:          "base matches release pattern, head matches feature",
			headRef:       "feature/abc",
			baseRef:       "release/v1.0",
			// multi-both-or matches because head=feature/abc matches HeadBranch patterns
			expectMatched: []string{"multi-head-or", "multi-base-or", "multi-both-or"},
		},
		{
			name:          "no head match, but base matches main",
			headRef:       "chore/cleanup",
			baseRef:       "main",
			// multi-both-or matches because base=main matches BaseBranch patterns
			expectMatched: []string{"multi-base-or", "multi-both-or"},
		},
		{
			name:          "head matches feature, but no base match",
			headRef:       "feature/test",
			baseRef:       "staging",
			// multi-both-or matches because head=feature/test matches HeadBranch patterns
			expectMatched: []string{"multi-head-or", "multi-both-or"},
		},
		{
			name:          "neither head nor base matches for multi-both-or",
			headRef:       "chore/cleanup",
			baseRef:       "staging",
			expectMatched: []string{}, // None should match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &github.PullRequest{
				Base:   &github.PullRequestBranch{Ref: github.Ptr(tt.baseRef)},
				Head:   &github.PullRequestBranch{Ref: github.Ptr(tt.headRef)},
				Labels: []*github.Label{},
			}
			files := []*github.CommitFile{}
			matcher := NewMatcher(context.TODO(), nil)
			result := matcher.CheckMatchConfigs(cfg, files, pr)

			for _, label := range tt.expectMatched {
				if !result.IsMatched(label) {
					t.Errorf("%s should be matched (head=%s, base=%s)", label, tt.headRef, tt.baseRef)
				}
			}
			// Check labels that should NOT be matched
			allLabels := []string{"multi-head-or", "multi-base-or", "multi-both-or"}
			for _, label := range allLabels {
				shouldMatch := false
				for _, expected := range tt.expectMatched {
					if label == expected {
						shouldMatch = true
						break
					}
				}
				if !shouldMatch && result.IsMatched(label) {
					t.Errorf("%s should NOT be matched (head=%s, base=%s)", label, tt.headRef, tt.baseRef)
				}
			}
		})
	}
}

// TestCheckMatchConfigs_MultipleBranchPatternsInAll tests that multiple branch patterns
// within a single LabelerRule under "all" are evaluated as AND (all must match)
func TestCheckMatchConfigs_MultipleBranchPatternsInAll(t *testing.T) {
	cfg := LabelerConfig{
		"all-head-and-base": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{All: []LabelerRule{{
					HeadBranch: []any{"feature/.*"},
					BaseBranch: []any{"main"},
				}}},
			},
		},
	}

	tests := []struct {
		name        string
		headRef     string
		baseRef     string
		shouldMatch bool
	}{
		{
			name:        "both match",
			headRef:     "feature/test",
			baseRef:     "main",
			shouldMatch: true,
		},
		{
			name:        "only head matches",
			headRef:     "feature/test",
			baseRef:     "develop",
			shouldMatch: false,
		},
		{
			name:        "only base matches",
			headRef:     "bugfix/test",
			baseRef:     "main",
			shouldMatch: false,
		},
		{
			name:        "neither matches",
			headRef:     "bugfix/test",
			baseRef:     "develop",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &github.PullRequest{
				Base:   &github.PullRequestBranch{Ref: github.Ptr(tt.baseRef)},
				Head:   &github.PullRequestBranch{Ref: github.Ptr(tt.headRef)},
				Labels: []*github.Label{},
			}
			files := []*github.CommitFile{}
			matcher := NewMatcher(context.TODO(), nil)
			result := matcher.CheckMatchConfigs(cfg, files, pr)

			if tt.shouldMatch && !result.IsMatched("all-head-and-base") {
				t.Errorf("all-head-and-base should be matched (head=%s, base=%s)", tt.headRef, tt.baseRef)
			}
			if !tt.shouldMatch && result.IsMatched("all-head-and-base") {
				t.Errorf("all-head-and-base should NOT be matched (head=%s, base=%s)", tt.headRef, tt.baseRef)
			}
		})
	}
}

// TestCheckMatchConfigs_MultipleRulesInAny tests that multiple LabelerRule entries
// in the "any" array are evaluated as OR (any rule succeeding is enough)
func TestCheckMatchConfigs_MultipleRulesInAny(t *testing.T) {
	cfg := LabelerConfig{
		"separate-rules-or": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{
					{HeadBranch: []any{"feature/.*"}},
					{HeadBranch: []any{"bugfix/.*"}},
					{BaseBranch: []any{"main"}},
				}},
			},
		},
	}

	tests := []struct {
		name        string
		headRef     string
		baseRef     string
		shouldMatch bool
	}{
		{
			name:        "first rule matches",
			headRef:     "feature/test",
			baseRef:     "develop",
			shouldMatch: true,
		},
		{
			name:        "second rule matches",
			headRef:     "bugfix/fix",
			baseRef:     "develop",
			shouldMatch: true,
		},
		{
			name:        "third rule matches",
			headRef:     "chore/cleanup",
			baseRef:     "main",
			shouldMatch: true,
		},
		{
			name:        "no rule matches",
			headRef:     "chore/cleanup",
			baseRef:     "develop",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &github.PullRequest{
				Base:   &github.PullRequestBranch{Ref: github.Ptr(tt.baseRef)},
				Head:   &github.PullRequestBranch{Ref: github.Ptr(tt.headRef)},
				Labels: []*github.Label{},
			}
			files := []*github.CommitFile{}
			matcher := NewMatcher(context.TODO(), nil)
			result := matcher.CheckMatchConfigs(cfg, files, pr)

			if tt.shouldMatch && !result.IsMatched("separate-rules-or") {
				t.Errorf("separate-rules-or should be matched (head=%s, base=%s)", tt.headRef, tt.baseRef)
			}
			if !tt.shouldMatch && result.IsMatched("separate-rules-or") {
				t.Errorf("separate-rules-or should NOT be matched (head=%s, base=%s)", tt.headRef, tt.baseRef)
			}
		})
	}
}
