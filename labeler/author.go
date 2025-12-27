package labeler

import (
	"context"
	"strings"

	"github.com/google/go-github/v79/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

// AuthorMatcher handles author matching logic including team membership checks
type AuthorMatcher struct {
	ctx context.Context
	g   *gh.GitHubClient
	// teamMembershipCache caches team membership results to avoid repeated API calls
	teamMembershipCache map[string]bool
}

// NewAuthorMatcher creates a new AuthorMatcher instance
func NewAuthorMatcher(ctx context.Context, g *gh.GitHubClient) *AuthorMatcher {
	return &AuthorMatcher{
		ctx:                 ctx,
		g:                   g,
		teamMembershipCache: make(map[string]bool),
	}
}

// MatchAuthor checks if the PR author matches any of the given patterns
// Patterns can be:
// - Regular expressions (matched against author username)
// - @org/team-slug (matches if author is a member of the team)
// - !@org/team-slug (matches if author is NOT a member of the team)
func (m *AuthorMatcher) MatchAuthor(patterns []string, pr *github.PullRequest) bool {
	if len(patterns) == 0 {
		return false
	}

	author := pr.GetUser().GetLogin()
	if author == "" {
		return false
	}

	for _, pattern := range patterns {
		if m.matchPattern(pattern, author) {
			return true
		}
	}
	return false
}

// matchPattern matches a single pattern against the author
func (m *AuthorMatcher) matchPattern(pattern, author string) bool {
	// Check for negated team pattern: !@org/team-slug
	if strings.HasPrefix(pattern, "!@") && strings.Contains(pattern[2:], "/") {
		// If no client is available, skip team patterns
		if m.g == nil || m.ctx == nil {
			return false
		}
		return m.matchNegatedTeam(pattern[2:], author)
	}

	// Check for team pattern: @org/team-slug
	if strings.HasPrefix(pattern, "@") && strings.Contains(pattern[1:], "/") {
		// If no client is available, skip team patterns
		if m.g == nil || m.ctx == nil {
			return false
		}
		return m.matchTeam(pattern[1:], author)
	}

	// Otherwise, treat as regex pattern
	return matchAnyRegex([]string{pattern}, author)
}

// matchTeam checks if the author is a member of the specified team
func (m *AuthorMatcher) matchTeam(teamRef, author string) bool {
	parts := strings.SplitN(teamRef, "/", 2)
	if len(parts) != 2 {
		return false
	}
	org := parts[0]
	teamSlug := parts[1]

	cacheKey := teamRef + ":" + author
	if cached, ok := m.teamMembershipCache[cacheKey]; ok {
		return cached
	}

	isMember := m.checkTeamMembership(org, teamSlug, author)
	m.teamMembershipCache[cacheKey] = isMember
	return isMember
}

// matchNegatedTeam checks if the author is NOT a member of the specified team
func (m *AuthorMatcher) matchNegatedTeam(teamRef, author string) bool {
	return !m.matchTeam(teamRef, author)
}

// checkTeamMembership checks if the user is a member of the specified team
func (m *AuthorMatcher) checkTeamMembership(org, teamSlug, username string) bool {
	if m.g == nil || m.ctx == nil {
		return false
	}

	// Use FindTeamMembership which returns nil without error if not found
	membership, err := m.g.FindTeamMembership(m.ctx, org, teamSlug, username)
	if err != nil {
		// If we can't check membership (e.g., API error, no permission), assume not a member
		return false
	}
	// membership.State can be "active" or "pending"
	// We consider both as a member for matching purposes
	return membership != nil && membership.GetState() != ""
}
