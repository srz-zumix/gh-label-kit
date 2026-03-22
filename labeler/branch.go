package labeler

import (
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

func matchLabelerRuleBaseBranch(r LabelerRule, pr *PullRequest) bool {
	if base := r.GetBaseBranch(); len(base) > 0 {
		for _, re := range base {
			if matchAnyRegex([]string{re}, pr.Base.GetRef()) {
				logger.Debug("BaseBranch pattern matched", "pattern", re, "branch", pr.Base.GetRef())
				return true
			}
		}
		logger.Debug("BaseBranch pattern not matched", "patterns", base, "branch", pr.Base.GetRef())
	}
	return false
}

func matchLabelerRuleHeadBranch(r LabelerRule, pr *PullRequest) bool {
	if head := r.GetHeadBranch(); len(head) > 0 {
		for _, re := range head {
			if matchAnyRegex([]string{re}, pr.Head.GetRef()) {
				logger.Debug("HeadBranch pattern matched", "pattern", re, "branch", pr.Head.GetRef())
				return true
			}
		}
		logger.Debug("HeadBranch pattern not matched", "patterns", head, "branch", pr.Head.GetRef())
	}
	return false
}
