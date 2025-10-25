package labeler

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/google/go-github/v73/github"
	"github.com/srz-zumix/go-gh-extension/pkg/actions"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

const (
	ReviewRequestModeNone             string = "none"
	ReviewRequestModeAddTo            string = "addto"
	ReviewRequestModeAlways           string = "always"
	ReviewRequestModeReadyForReview   string = "ready_for_review"
	ReviewRequestModeAlwaysReviewable string = "always_reviewable"
)

var ReviewersRequestModes = []string{
	ReviewRequestModeNone,
	ReviewRequestModeAddTo,
	ReviewRequestModeAlways,
	ReviewRequestModeReadyForReview,
	ReviewRequestModeAlwaysReviewable,
}

func GetReviewRequestTargetLabels(pr *github.PullRequest, matchResult MatchResult, reviewRequestMode string, syncLabels bool) []string {
	switch reviewRequestMode {
	case ReviewRequestModeNone:
		return nil
	case ReviewRequestModeAddTo:
		return matchResult.AddTo()
	case ReviewRequestModeAlways:
		return matchResult.GetLabels(syncLabels)
	case ReviewRequestModeReadyForReview:
		if pr.GetDraft() {
			return nil
		}
		if actions.IsRunsOn() {
			eventContext, err := actions.GetEventPayload()
			if err != nil {
				return nil
			}
			if eventContext.Action != "ready_for_review" {
				return nil
			}
			return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAlways, syncLabels)
		}
		return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAddTo, syncLabels)
	case ReviewRequestModeAlwaysReviewable:
		if pr.GetDraft() {
			return nil
		}
		return GetReviewRequestTargetLabels(pr, matchResult, ReviewRequestModeAlways, syncLabels)
	default:
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
	if len(labels) == 0 {
		return []string{}
	}
	if c.pr.GetState() != "open" {
		return []string{}
	}

	codeowners := CollectCodeownersSet(labels, c.cfg)
	codeowners = c.ExpandCodeownersSet(codeowners)
	author := c.pr.GetUser().GetLogin()
	delete(codeowners, author)
	reviewers, err := gh.ListPullRequestReviewers(c.ctx, c.g, c.repo, c.pr)
	if err == nil {
		for _, r := range reviewers.Users {
			delete(codeowners, r.GetLogin())
		}
	}
	reviewed_reviewers, err := gh.GetPullRequestLatestReviews(c.ctx, c.g, c.repo, c.pr)
	if err == nil {
		for _, r := range reviewed_reviewers {
			if r.GetCommitID() == c.pr.GetHead().GetSHA() {
				delete(codeowners, r.GetUser().GetLogin())
			}
		}
	}
	return slices.Collect(maps.Keys(codeowners))
}

func (c *LabeledCodeOwners) SetReviewers(labels []string) ([]string, *github.PullRequest, error) {
	codeowners := c.GetReviewers(labels)
	if len(codeowners) == 0 {
		return nil, c.pr, nil
	}
	pr, err := gh.RequestPullRequestReviewers(c.ctx, c.g, c.repo, c.pr, gh.GetRequestedReviewers(codeowners))
	return codeowners, pr, err
}
