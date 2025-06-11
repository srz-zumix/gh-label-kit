package labeler

import (
	"context"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v71/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

func SetLabels(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, pr *github.PullRequest, matchResult MatchResult, sync bool) ([]*github.Label, error) {
	allLebels := matchResult.GetLabels(sync)
	return gh.SetPullRequestLabelsByNumber(ctx, g, repo, pr.GetNumber(), allLebels)
}
