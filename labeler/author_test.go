package labeler

import (
	"context"
	"testing"

	"github.com/google/go-github/v79/github"
)

func TestMatchLabelerRuleAuthor_RegexMatch(t *testing.T) {
	tests := []struct {
		name       string
		patterns   []string
		authorName string
		want       bool
	}{
		{
			name:       "exact match",
			patterns:   []string{"user1"},
			authorName: "user1",
			want:       true,
		},
		{
			name:       "regex pattern match",
			patterns:   []string{"user.*"},
			authorName: "user123",
			want:       true,
		},
		{
			name:       "regex pattern no match",
			patterns:   []string{"admin.*"},
			authorName: "user123",
			want:       false,
		},
		{
			name:       "multiple patterns - one matches",
			patterns:   []string{"admin.*", "user.*"},
			authorName: "user123",
			want:       true,
		},
		{
			name:       "multiple patterns - none match",
			patterns:   []string{"admin.*", "moderator.*"},
			authorName: "user123",
			want:       false,
		},
		{
			name:       "empty patterns",
			patterns:   []string{},
			authorName: "user123",
			want:       false,
		},
		{
			name:       "empty author",
			patterns:   []string{"user.*"},
			authorName: "",
			want:       false,
		},
		{
			name:       "case sensitive regex",
			patterns:   []string{"User.*"},
			authorName: "user123",
			want:       false,
		},
		{
			name:       "case insensitive regex with flag",
			patterns:   []string{"(?i)User.*"},
			authorName: "user123",
			want:       true,
		},
		{
			name:       "bot user pattern",
			patterns:   []string{".*\\[bot\\]$"},
			authorName: "dependabot[bot]",
			want:       true,
		},
		{
			name:       "prefix match",
			patterns:   []string{"^feature-"},
			authorName: "feature-user",
			want:       true,
		},
		{
			name:       "suffix match",
			patterns:   []string{"-bot$"},
			authorName: "github-bot",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := LabelerRule{
				Author: tt.patterns,
			}
			pr := &github.PullRequest{
				User: &github.User{Login: github.Ptr(tt.authorName)},
			}
			// Use Matcher with context.TODO() and nil client to test regex-only matching
			matcher := NewMatcher(context.TODO(), nil)
			got := matcher.matchLabelerRuleAuthor(rule, pr)
			if got != tt.want {
				t.Errorf("matchLabelerRuleAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorMatcher_MatchAuthor_RegexOnly(t *testing.T) {
	// Test AuthorMatcher without API client (regex patterns only)
	matcher := NewAuthorMatcher(context.TODO(), nil)

	tests := []struct {
		name       string
		patterns   []string
		authorName string
		want       bool
	}{
		{
			name:       "exact match",
			patterns:   []string{"testuser"},
			authorName: "testuser",
			want:       true,
		},
		{
			name:       "regex match",
			patterns:   []string{"test.*"},
			authorName: "testuser",
			want:       true,
		},
		{
			name:       "no match",
			patterns:   []string{"admin"},
			authorName: "testuser",
			want:       false,
		},
		{
			name:       "team pattern without client - should skip",
			patterns:   []string{"@org/team"},
			authorName: "testuser",
			want:       false,
		},
		{
			name:       "negated team pattern without client - should skip",
			patterns:   []string{"!@org/team"},
			authorName: "testuser",
			want:       false,
		},
		{
			name:       "mixed patterns - regex matches",
			patterns:   []string{"@org/team", "test.*"},
			authorName: "testuser",
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &github.PullRequest{
				User: &github.User{Login: github.Ptr(tt.authorName)},
			}
			got := matcher.MatchAuthor(tt.patterns, pr)
			if got != tt.want {
				t.Errorf("AuthorMatcher.MatchAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorMatcher_matchPattern(t *testing.T) {
	matcher := NewAuthorMatcher(context.TODO(), nil)

	tests := []struct {
		name    string
		pattern string
		author  string
		want    bool
	}{
		{
			name:    "simple regex",
			pattern: "user.*",
			author:  "user123",
			want:    true,
		},
		{
			name:    "team pattern - skipped without client",
			pattern: "@org/team",
			author:  "user123",
			want:    false,
		},
		{
			name:    "negated team pattern - skipped without client",
			pattern: "!@org/team",
			author:  "user123",
			want:    false,
		},
		{
			name:    "@ without slash is regex",
			pattern: "@user",
			author:  "@user",
			want:    true,
		},
		{
			name:    "!@ without slash is regex",
			pattern: "!@user",
			author:  "!@user",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.matchPattern(tt.pattern, tt.author)
			if got != tt.want {
				t.Errorf("AuthorMatcher.matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.author, got, tt.want)
			}
		})
	}
}

func TestCheckMatchConfigs_Author(t *testing.T) {
	cfg := LabelerConfig{
		"bot-label": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{Author: []any{".*\\[bot\\]$"}}}},
			},
		},
		"user-label": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{Any: []LabelerRule{{Author: []any{"testuser"}}}},
			},
		},
		"combined-label": LabelerLabelConfig{
			Matcher: []LabelerMatch{
				{All: []LabelerRule{
					{Author: []any{"testuser"}},
					{HeadBranch: []any{"feature/.*"}},
				}},
			},
		},
	}

	tests := []struct {
		name         string
		authorName   string
		headBranch   string
		wantMatched  []string
		wantUnmatched []string
	}{
		{
			name:         "bot user matches bot-label",
			authorName:   "dependabot[bot]",
			headBranch:   "main",
			wantMatched:  []string{"bot-label"},
			wantUnmatched: []string{"combined-label", "user-label"},
		},
		{
			name:         "testuser matches user-label",
			authorName:   "testuser",
			headBranch:   "main",
			wantMatched:  []string{"user-label"},
			wantUnmatched: []string{"bot-label", "combined-label"},
		},
		{
			name:         "testuser with feature branch matches combined-label and user-label",
			authorName:   "testuser",
			headBranch:   "feature/test",
			wantMatched:  []string{"combined-label", "user-label"},
			wantUnmatched: []string{"bot-label"},
		},
		{
			name:         "other user matches nothing",
			authorName:   "otheruser",
			headBranch:   "main",
			wantMatched:  []string{},
			wantUnmatched: []string{"bot-label", "combined-label", "user-label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &github.PullRequest{
				User:   &github.User{Login: github.Ptr(tt.authorName)},
				Head:   &github.PullRequestBranch{Ref: github.Ptr(tt.headBranch)},
				Base:   &github.PullRequestBranch{Ref: github.Ptr("main")},
				Labels: []*github.Label{},
			}
			files := []*github.CommitFile{}
			matcher := NewMatcher(context.TODO(), nil)
			result := matcher.CheckMatchConfigs(cfg, files, pr)

			for _, want := range tt.wantMatched {
				if !result.IsMatched(want) {
					t.Errorf("%s should be matched", want)
				}
			}
			for _, want := range tt.wantUnmatched {
				if !result.IsUnmatched(want) {
					t.Errorf("%s should be unmatched", want)
				}
			}
		})
	}
}
