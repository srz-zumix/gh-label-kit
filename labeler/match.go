package labeler

import (
	"maps"
	"path/filepath"
	"slices"

	"github.com/dlclark/regexp2"
	"github.com/google/go-github/v71/github"
)

type MatchResult struct {
	Current   []string // Current labels on the PR
	Matched   []string // Matched label names
	Unmatched []string // Unmatched label names
}

func (r MatchResult) GetLabels(sync bool) []string {
	if sync {
		return r.SyncTo()
	}
	return r.SetTo()
}

func (r MatchResult) SetTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Current {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	return slices.Collect(maps.Keys(allLabels))
}

func (r MatchResult) SyncTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Current {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Unmatched {
		delete(allLabels, label)
	}
	return slices.Collect(maps.Keys(allLabels))
}

func (r MatchResult) AddTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Matched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Current {
		delete(allLabels, label)
	}
	return slices.Collect(maps.Keys(allLabels))
}

func (r MatchResult) DeleteTo() []string {
	allLabels := make(map[string]struct{})
	for _, label := range r.Unmatched {
		allLabels[label] = struct{}{}
	}
	for _, label := range r.Current {
		if _, exists := allLabels[label]; !exists {
			delete(allLabels, label)
		}
	}
	return slices.Collect(maps.Keys(allLabels))
}

func CheckMatchConfigs(cfg LabelerConfig, changedFiles []*github.CommitFile, pr *github.PullRequest) MatchResult {
	result := MatchResult{
		Current:   []string{},
		Matched:   []string{},
		Unmatched: []string{},
	}

	for _, label := range pr.Labels {
		result.Current = append(result.Current, label.GetName())
	}

	for label, matches := range cfg {
		matched := false
		for _, match := range matches {
			if MatchLabelerMatch(match, changedFiles, pr) {
				matched = true
				break
			}
		}
		if matched {
			result.Matched = append(result.Matched, label)
		} else {
			result.Unmatched = append(result.Unmatched, label)
		}
	}

	return result
}

// matchLabelerMatch checks if a PR matches a label's match object (any/all/changed-files/branch)
func MatchLabelerMatch(m LabelerMatch, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	if len(m.Any) > 0 {
		for _, rule := range m.Any {
			if MatchLabelerRule(rule, changedFiles, pr) {
				return true
			}
		}
		return false
	}
	if len(m.All) > 0 {
		for _, rule := range m.All {
			if !MatchLabelerRule(rule, changedFiles, pr) {
				return false
			}
		}
		return true
	}
	// ベース互換: any/all省略時はany扱い
	return MatchLabelerRule(LabelerRule{
		ChangedFiles: m.ChangedFiles,
		BaseBranch:   m.BaseBranch,
		HeadBranch:   m.HeadBranch,
	}, changedFiles, pr)
}

func MatchLabelerRule(r LabelerRule, changedFiles []*github.CommitFile, pr *github.PullRequest) bool {
	if len(r.ChangedFiles) > 0 {
		for _, cf := range r.ChangedFiles {
			if MatchChangedFilesRule(cf, changedFiles) {
				return true
			}
		}
	}
	if base := r.GetBaseBranch(); len(base) > 0 {
		for _, re := range base {
			if MatchAnyRegex([]string{re}, pr.Base.GetRef()) {
				return true
			}
		}
	}
	if head := r.GetHeadBranch(); len(head) > 0 {
		for _, re := range head {
			if MatchAnyRegex([]string{re}, pr.Head.GetRef()) {
				return true
			}
		}
	}
	return false
}

func MatchChangedFilesRule(cf ChangedFilesRule, changedFiles []*github.CommitFile) bool {
	// any-glob-to-any-file: ANY glob matches ANY file
	for _, pattern := range cf.AnyGlobToAnyFile {
		for _, f := range changedFiles {
			if f.Filename != nil && matchGlob(pattern, *f.Filename) {
				return true
			}
		}
	}
	// any-glob-to-all-files: ANY glob matches ALL files
	if len(cf.AnyGlobToAllFiles) > 0 && len(changedFiles) > 0 {
		for _, pattern := range cf.AnyGlobToAllFiles {
			allMatch := true
			for _, f := range changedFiles {
				if f.Filename == nil || !matchGlob(pattern, *f.Filename) {
					allMatch = false
					break
				}
			}
			if allMatch {
				return true
			}
		}
	}
	// all-globs-to-any-file: ALL globs match at least one file
	if len(cf.AllGlobsToAnyFile) > 0 {
		for _, pattern := range cf.AllGlobsToAnyFile {
			found := false
			for _, f := range changedFiles {
				if f.Filename != nil && matchGlob(pattern, *f.Filename) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	// all-globs-to-all-files: ALL globs match ALL files
	if len(cf.AllGlobsToAllFiles) > 0 && len(changedFiles) > 0 {
		for _, pattern := range cf.AllGlobsToAllFiles {
			for _, f := range changedFiles {
				if f.Filename == nil || !matchGlob(pattern, *f.Filename) {
					return false
				}
			}
		}
		return true
	}
	return false
}

func matchGlob(pattern, filename string) bool {
	matched, err := filepath.Match(pattern, filename)
	return err == nil && matched
}

func MatchAnyRegex(patterns []string, branch string) bool {
	for _, pattern := range patterns {
		re := regexp2.MustCompile(pattern, regexp2.RE2)
		matched, err := re.MatchString(branch)
		if err == nil && matched {
			return true
		}
	}
	return false
}
