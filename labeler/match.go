package labeler

import (
	"context"
	"maps"
	"slices"

	"github.com/google/go-github/v79/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

// Matcher handles all matching logic for labeler config rules
type Matcher struct {
	authorMatcher *AuthorMatcher
}

// NewMatcher creates a new Matcher instance with the given context and GitHub client
func NewMatcher(ctx context.Context, g *gh.GitHubClient) *Matcher {
	return &Matcher{
		authorMatcher: NewAuthorMatcher(ctx, g),
	}
}

type MatchResult struct {
	Current   []string // Current labels on the PR
	Matched   []string // Matched label names
	Unmatched []string // Unmatched label names
}

func (r MatchResult) GetLabels(sync bool) []string {
	if sync {
		return r.SyncTo()
	}
	return r.SetTo()
}

func (r MatchResult) IsMatched(label string) bool {
	for _, matched := range r.Matched {
		if matched == label {
			return true
		}
	}
	return false
}

func (r MatchResult) IsUnmatched(label string) bool {
	for _, matched := range r.Unmatched {
		if matched == label {
			return true
		}
	}
	return false
}

func (r MatchResult) HasDiff(sync bool) bool {
	allLabels := make(map[string]struct{})
	for _, label := range r.Current {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Matched {
		if _, exists := allLabels[label]; !exists {
			return true
		}
	}
	if sync {
		for _, label := range r.Unmatched {
			if _, exists := allLabels[label]; exists {
				return true
			}
		}
	}
	return false

}

func (r MatchResult) SetTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Current {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	labels := slices.Collect(maps.Keys(allLabels))
	slices.Sort(labels)
	return labels
}

func (r MatchResult) SyncTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Current {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Unmatched {
		delete(allLabels, label)
	}
	labels := slices.Collect(maps.Keys(allLabels))
	slices.Sort(labels)
	return labels
}

func (r MatchResult) AddTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Current {
		delete(allLabels, label)
	}
	labels := slices.Collect(maps.Keys(allLabels))
	slices.Sort(labels)
	return labels
}

func (r MatchResult) DeleteTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Unmatched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Current {
		if _, exists := allLabels[label]; !exists {
			delete(allLabels, label)
		}
	}
	labels := slices.Collect(maps.Keys(allLabels))
	slices.Sort(labels)
	return labels
}

// CheckMatchConfigs checks all label configs against the PR and returns matched/unmatched labels
func (m *Matcher) CheckMatchConfigs(cfg LabelerConfig, changedFiles []*github.CommitFile, pr *github.PullRequest) MatchResult {
	result := MatchResult{
		Current:   []string{},
		Matched:   []string{},
		Unmatched: []string{},
	}

	for _, label := range pr.Labels {
		result.Current = append(result.Current, label.GetName())
	}

	for label, labelConfig := range cfg {
		matched := len(labelConfig.Matcher) != 0
		for _, match := range labelConfig.Matcher {
			if !m.matchLabelerMatch(match, changedFiles, pr) {
				matched = false
				break
			}
		}
		if matched {
			result.Matched = append(result.Matched, label)
		} else {
			result.Unmatched = append(result.Unmatched, label)
		}
	}
	slices.Sort(result.Current)
	slices.Sort(result.Matched)
	slices.Sort(result.Unmatched)
	return result
}

// matchLabelerMatch checks if a PR matches a label's match object (any/all/changed-files/branch/author)
func (m *Matcher) matchLabelerMatch(match LabelerMatch, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	if len(match.All) > 0 {
		if !m.matchLabelerMatchAll(match.All, changedFiles, pr) {
			return false
		}
	}
	if len(match.Any) > 0 {
		if !m.matchLabelerMatchAny(match.Any, changedFiles, pr) {
			return false
		}
	}
	return true
}

func (m *Matcher) matchLabelerMatchAny(rules []LabelerRule, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	for _, rule := range rules {
		if m.matchLabelerRuleAny(rule, changedFiles, pr) {
			return true
		}
	}
	return false
}

func (m *Matcher) matchLabelerMatchAll(rules []LabelerRule, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	for _, rule := range rules {
		if !m.matchLabelerRuleAll(rule, changedFiles, pr) {
			return false
		}
	}
	return true
}

func (m *Matcher) matchLabelerRuleAny(r LabelerRule, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	if r.BaseBranch != nil {
		if matchLabelerRuleBaseBranch(r, pr) {
			return true
		}
	}
	if r.HeadBranch != nil {
		if matchLabelerRuleHeadBranch(r, pr) {
			return true
		}
	}
	if r.Author != nil {
		if m.matchLabelerRuleAuthor(r, pr) {
			return true
		}
	}
	if len(r.ChangedFiles) > 0 {
		if matchChangedFilesAny(r.ChangedFiles, changedFiles) {
			return true
		}
	}
	return false
}

func (m *Matcher) matchLabelerRuleAll(r LabelerRule, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	if r.BaseBranch != nil {
		if !matchLabelerRuleBaseBranch(r, pr) {
			return false
		}
	}
	if r.HeadBranch != nil {
		if !matchLabelerRuleHeadBranch(r, pr) {
			return false
		}
	}
	if r.Author != nil {
		if !m.matchLabelerRuleAuthor(r, pr) {
			return false
		}
	}
	if len(r.ChangedFiles) > 0 {
		if !matchChangedFilesAll(r.ChangedFiles, changedFiles) {
			return false
		}
	}
	return true
}

// matchLabelerRuleAuthor checks if the PR author matches the rule's author patterns
func (m *Matcher) matchLabelerRuleAuthor(r LabelerRule, pr *github.PullRequest) bool {
	authors := r.GetAuthor()
	if len(authors) == 0 {
		return false
	}
	return m.authorMatcher.MatchAuthor(authors, pr)
}
