package labeler

import (
	"github.com/dlclark/regexp2"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

func matchAnyRegex(patterns []string, str string) bool {
	for _, pattern := range patterns {
		re := regexp2.MustCompile(pattern, regexp2.RE2)
		matched, err := re.MatchString(str)
		if err == nil && matched {
			logger.Debug("Regex pattern matched", "pattern", pattern, "string", str)
			return true
		}
	}
	logger.Debug("No regex pattern matched", "patterns", patterns, "string", str)
	return false
}
