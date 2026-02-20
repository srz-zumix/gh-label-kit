package labeler

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v79/github"
	"github.com/srz-zumix/go-gh-extension/pkg/actions"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

const (
	ReviewRequestModeNone             string = "none"
	ReviewRequestModeNever            string = "never"
	ReviewRequestModeAddTo            string = "addto"
	ReviewRequestModeAlways           string = "always"
	ReviewRequestModeReadyForReview   string = "ready_for_review"
	ReviewRequestModeAlwaysReviewable string = "always_reviewable"
)

var ReviewersRequestModes = []string{
	ReviewRequestModeAddTo,
	ReviewRequestModeAlways,
	ReviewRequestModeReadyForReview,
	ReviewRequestModeAlwaysReviewable,
	ReviewRequestModeNever,
	ReviewRequestModeNone,
}

func GetReviewRequestTargetLabels(pr *github.PullRequest, matchResult MatchResult, reviewRequestMode string, syncLabels bool) []string {
	logger.Debug("Getting review request target labels", "pr", pr.GetNumber(), "mode", reviewRequestMode, "syncLabels", syncLabels)
	switch reviewRequestMode {
	case ReviewRequestModeNone, ReviewRequestModeNever:
		logger.Debug("Review request mode is none/never", "mode", reviewRequestMode)
		return nil
	case ReviewRequestModeAddTo:
		labels := matchResult.AddTo()
		logger.Debug("Review request mode addto", "labels", labels)
		return labels
	case ReviewRequestModeAlways:
		labels := matchResult.GetLabels(syncLabels)
		logger.Debug("Review request mode always", "labels", labels)
		return labels
	case ReviewRequestModeReadyForReview:
		if pr.GetDraft() {
			logger.Debug("PR is draft, skipping review request", "pr", pr.GetNumber())
			return nil
		}
		if actions.IsRunsOn() {
			eventContext, err := actions.GetEventPayload()
			if err != nil {
				return nil
			}
			if eventContext.Action == "ready_for_review" {
				logger.Debug("Event action is ready_for_review, using always mode", "pr", pr.GetNumber())
				return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAlways, syncLabels)
			}
		}
		return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAddTo, syncLabels)
	case ReviewRequestModeAlwaysReviewable:
		if pr.GetDraft() {
			logger.Debug("PR is draft, skipping review request", "pr", pr.GetNumber())
			return nil
		}
		return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAlways, syncLabels)
	default:
		logger.Debug("Unknown review request mode", "mode", reviewRequestMode)
		return nil
	}
}

func CollectCodeownersSet(labels []string, cfg LabelerConfig) map[string]struct{} {
	ownerSet := make(map[string]struct{})
	for _, label := range labels {
		if lc, ok := cfg[label]; ok {
			for _, owner := range lc.Codeowners {
				if owner[0] == '@' {
					owner = owner[1:]
				}
				ownerSet[owner] = struct{}{}
			}
		}
	}
	return ownerSet
}

func CollectCodeowners(labels []string, cfg LabelerConfig) []string {
	ownerSet := CollectCodeownersSet(labels, cfg)
	return slices.Collect(maps.Keys(ownerSet))
}

type LabeledCodeOwners struct {
	ctx  context.Context
	g    *gh.GitHubClient
	repo repository.Repository
	cfg  LabelerConfig
	pr   *github.PullRequest
	mode string
}

func NewLabeledCodeOwners(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, pr *github.PullRequest, cfg LabelerConfig, mode string) *LabeledCodeOwners {
	return &LabeledCodeOwners{
		ctx:  ctx,
		g:    g,
		repo: repo,
		cfg:  cfg,
		pr:   pr,
		mode: mode,
	}
}

func (c *LabeledCodeOwners) ExpandCodeownersSet(ownerSet map[string]struct{}) map[string]struct{} {
	expandedOwnerSet := make(map[string]struct{})
	for config := range ownerSet {
		if len(config) >= 3 && config[len(config)-3:] == "..." {
			owner := config[:len(config)-3]
			if strings.Contains(owner, "/") {
				parts := strings.SplitN(owner, "/", 2)
				members, err := gh.ListTeamMembers(c.ctx, c.g, c.repo, parts[1], nil, false)
				if err == nil {
					for _, member := range members {
						expandedOwnerSet[member.GetLogin()] = struct{}{}
					}
					continue
				} else {
					fmt.Fprintf(os.Stderr, "failed to expand team %s: %v", owner, err)
				}
			}
			expandedOwnerSet[owner] = struct{}{}
		} else {
			expandedOwnerSet[config] = struct{}{}
		}
	}
	return expandedOwnerSet
}

func (c *LabeledCodeOwners) GetReviewers(labels []string) []string {
	logger.Debug("Getting reviewers from labels", "labelsCount", len(labels))
	if len(labels) == 0 {
		return []string{}
	}
	if c.pr.GetState() != "open" {
		logger.Debug("PR is not open, skipping reviewers", "pr", c.pr.GetNumber(), "state", c.pr.GetState())
		return []string{}
	}

	codeowners := CollectCodeownersSet(labels, c.cfg)
	logger.Debug("Collected codeowners", "count", len(codeowners))
	codeowners = c.ExpandCodeownersSet(codeowners)
	logger.Debug("Expanded codeowners", "count", len(codeowners))
	author := c.pr.GetUser().GetLogin()
	delete(codeowners, author)
	logger.Debug("Removed PR author from codeowners", "author", author)
	reviewers, err := gh.ListPullRequestReviewers(c.ctx, c.g, c.repo, c.pr)
	if err == nil {
		for _, r := range reviewers.Users {
			delete(codeowners, r.GetLogin())
		}
		logger.Debug("Removed existing reviewers", "reviewersCount", len(reviewers.Users))
	}
	reviewedReviewers, err := gh.GetPullRequestLatestReviews(c.ctx, c.g, c.repo, c.pr)
	if err == nil {
		removedCount := 0
		for _, r := range reviewedReviewers {
			if r.GetCommitID() == c.pr.GetHead().GetSHA() {
				delete(codeowners, r.GetUser().GetLogin())
				removedCount++
			}
		}
		logger.Debug("Removed reviewers who already reviewed latest commit", "removedCount", removedCount)
	}
	result := slices.Collect(maps.Keys(codeowners))
	logger.Debug("Final reviewers list", "reviewers", result)
	return result
}

func (c *LabeledCodeOwners) SetReviewers(labels []string) ([]string, *github.PullRequest, error) {
	codeowners := c.GetReviewers(labels)
	if len(codeowners) == 0 {
		logger.Debug("No reviewers to request", "pr", c.pr.GetNumber())
		return nil, c.pr, nil
	}
	logger.Debug("Requesting reviewers for PR", "pr", c.pr.GetNumber(), "reviewers", codeowners)
	pr, err := gh.RequestPullRequestReviewers(c.ctx, c.g, c.repo, c.pr, gh.GetRequestedReviewers(codeowners))
	if err != nil {
		logger.Debug("Failed to request reviewers", "pr", c.pr.GetNumber(), "error", err)
	}
	return codeowners, pr, err
}
