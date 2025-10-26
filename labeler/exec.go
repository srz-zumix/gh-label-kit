package labeler

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v73/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

func SetLabels(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, pr *github.PullRequest, allLebels []string, cfg LabelerConfig) ([]*github.Label, error) {
	var excessLabels []string
	if len(allLebels) > 100 {
		excessLabels = allLebels[100:]
		allLebels = allLebels[:100]
	}
	labels, err := gh.SetPullRequestLabels(ctx, g, repo, pr, allLebels)
	if err != nil {
		return nil, err
	}
	_, err = EditLabelsByConfig(ctx, g, repo, labels, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit labels for PR #%d: %w", pr.GetNumber(), err)
	}
	if len(excessLabels) > 0 {
		return labels, fmt.Errorf("label limit for a PR exceeded: not applied to PR #%d: %v", pr.GetNumber(), excessLabels)
	}
	return labels, nil
}
