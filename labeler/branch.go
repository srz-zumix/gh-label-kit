package labeler

import (
	"github.com/google/go-github/v71/github"
)

func matchLabelerRuleBaseBranch(r LabelerRule, pr *github.PullRequest) bool {
	if base := r.GetBaseBranch(); len(base) > 0 {
		for _, re := range base {
			if matchAnyRegex([]string{re}, pr.Base.GetRef()) {
				return true
			}
		}
	}
	return false
}

func matchLabelerRuleHeadBranch(r LabelerRule, pr *github.PullRequest) bool {
	if head := r.GetHeadBranch(); len(head) > 0 {
		for _, re := range head {
			if matchAnyRegex([]string{re}, pr.Head.GetRef()) {
				return true
			}
		}
	}
	return false
}
