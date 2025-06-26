package labeler

import (
	"github.com/dlclark/regexp2"
)

func matchAnyRegex(patterns []string, str string) bool {
	for _, pattern := range patterns {
		re := regexp2.MustCompile(pattern, regexp2.RE2)
		matched, err := re.MatchString(str)
		if err == nil && matched {
			return true
		}
	}
	return false
}
