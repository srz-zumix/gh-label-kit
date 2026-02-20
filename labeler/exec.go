package labeler

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v79/github"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

func SetLabels(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, pr *github.PullRequest, allLebels []string, cfg LabelerConfig) ([]*github.Label, error) {
	logger.Debug("Setting labels for PR", "pr", pr.GetNumber(), "labels", allLebels, "count", len(allLebels))
	var excessLabels []string
	if len(allLebels) > 100 {
		excessLabels = allLebels[100:]
		allLebels = allLebels[:100]
		logger.Debug("Label count exceeds limit, truncating", "pr", pr.GetNumber(), "limit", 100, "excess", excessLabels)
	}
	labels, err := gh.SetPullRequestLabels(ctx, g, repo, pr, allLebels)
	if err != nil {
		logger.Debug("Failed to set PR labels", "pr", pr.GetNumber(), "error", err)
		return nil, err
	}
	_, err = EditLabelsByConfig(ctx, g, repo, labels, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to edit labels for PR #%d: %w", pr.GetNumber(), err)
	}
	if len(excessLabels) > 0 {
		return labels, fmt.Errorf("label limit for a PR exceeded: not applied to PR #%d: %v", pr.GetNumber(), excessLabels)
	}
	logger.Debug("Successfully set labels for PR", "pr", pr.GetNumber(), "labels", len(labels))
	return labels, nil
}
