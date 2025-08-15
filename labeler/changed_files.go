package labeler

import (
	"github.com/google/go-github/v73/github"
)

func matchChangedFilesAny(rules []ChangedFilesRule, changedFiles []*github.CommitFile) bool {
	// Check if any of the rules match the changed files
	for _, rule := range rules {
		if matchChangedFilesRuleAny(rule, changedFiles) {
			return true
		}
	}
	return false
}

func matchChangedFilesAll(rules []ChangedFilesRule, changedFiles []*github.CommitFile) bool {
	// Check if any of the rules match the changed files
	for _, rule := range rules {
		if !matchChangedFilesRuleAll(rule, changedFiles) {
			return false
		}
	}
	return true
}

func matchAnyGlobToAnyFile(patterns []string, changedFiles []*github.CommitFile) bool {
	for _, pattern := range patterns {
		for _, f := range changedFiles {
			if f.Filename != nil && matchGlob(pattern, *f.Filename) {
				return true
			}
		}
	}
	return false
}

func matchAnyGlobToAllFiles(patterns []string, changedFiles []*github.CommitFile) bool {
	if len(changedFiles) == 0 {
		return false
	}
	for _, pattern := range patterns {
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
	return false
}

func matchAllGlobsToAnyFile(patterns []string, changedFiles []*github.CommitFile) bool {
	for _, pattern := range patterns {
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

func matchAllGlobsToAllFiles(patterns []string, changedFiles []*github.CommitFile) bool {
	if len(changedFiles) == 0 {
		return false
	}
	for _, pattern := range patterns {
		for _, f := range changedFiles {
			if f.Filename == nil || !matchGlob(pattern, *f.Filename) {
				return false
			}
		}
	}
	return true
}

// MatchChangedFilesRule checks if changed files match the given ChangedFilesRule.
func matchChangedFilesRuleAny(cf ChangedFilesRule, changedFiles []*github.CommitFile) bool {
	if len(cf.AnyGlobToAnyFile) != 0 {
		if matchAnyGlobToAnyFile(cf.AnyGlobToAnyFile, changedFiles) {
			return true
		}
	}
	if len(cf.AnyGlobToAllFiles) != 0 {
		if matchAnyGlobToAllFiles(cf.AnyGlobToAllFiles, changedFiles) {
			return true
		}
	}
	if len(cf.AllGlobsToAnyFile) != 0 {
		if matchAllGlobsToAnyFile(cf.AllGlobsToAnyFile, changedFiles) {
			return true
		}
	}
	if len(cf.AllGlobsToAllFiles) != 0 {
		if matchAllGlobsToAllFiles(cf.AllGlobsToAllFiles, changedFiles) {
			return true
		}
	}
	return false
}

func matchChangedFilesRuleAll(cf ChangedFilesRule, changedFiles []*github.CommitFile) bool {
	if len(cf.AnyGlobToAnyFile) != 0 {
		if !matchAnyGlobToAnyFile(cf.AnyGlobToAnyFile, changedFiles) {
			return false
		}
	}
	if len(cf.AnyGlobToAllFiles) != 0 {
		if !matchAnyGlobToAllFiles(cf.AnyGlobToAllFiles, changedFiles) {
			return false
		}
	}
	if len(cf.AllGlobsToAnyFile) != 0 {
		if !matchAllGlobsToAnyFile(cf.AllGlobsToAnyFile, changedFiles) {
			return false
		}
	}
	if len(cf.AllGlobsToAllFiles) != 0 {
		if !matchAllGlobsToAllFiles(cf.AllGlobsToAllFiles, changedFiles) {
			return false
		}
	}
	return true
}
