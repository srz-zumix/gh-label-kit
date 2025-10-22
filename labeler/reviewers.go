package labeler

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
)

const (
	ReviewRequestModeNone   string = "none"
	ReviewRequestModeAddTo  string = "addto"
	ReviewRequestModeAlways string = "always"
)

var ReviewersRequestModes = []string{
	ReviewRequestModeNone,
	ReviewRequestModeAddTo,
	ReviewRequestModeAlways,
}

func GetReviewRequestTargetLabels(matchResult MatchResult, reviewRequestMode string, syncLabels bool) []string {
	switch reviewRequestMode {
	case ReviewRequestModeNone:
		return nil
	case ReviewRequestModeAddTo:
		return matchResult.AddTo()
	case ReviewRequestModeAlways:
		return matchResult.GetLabels(syncLabels)
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

func ExpandCodeownersSet(ctx context.Context, g *gh.GitHubClient, repo repository.Repository, ownerSet map[string]struct{}, cfg LabelerConfig) map[string]struct{} {
	expandedOwnerSet := make(map[string]struct{})
	for config := range ownerSet {
		if len(config) >= 3 && config[len(config)-3:] == "..." {
			owner := config[:len(config)-3]
			if strings.Contains(owner, "/") {
				parts := strings.SplitN(owner, "/", 2)
				members, err := gh.ListTeamMembers(ctx, g, repo, parts[1], nil, false)
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

func CollectCodeowners(labels []string, cfg LabelerConfig) []string {
	ownerSet := CollectCodeownersSet(labels, cfg)
	return slices.Collect(maps.Keys(ownerSet))
}
