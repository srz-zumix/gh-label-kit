package labeler

import (
	"testing"
)

func TestMatchAnyRegex(t *testing.T) {
	cases := []struct {
		patterns []string
		branch   string
		want     bool
	}{
		{[]string{"^feature"}, "feature/abc", true},
		{[]string{"^hotfix"}, "feature/abc", false},
		{[]string{"^f.*e$"}, "feature", true},
		// 否定先読み
		{[]string{"^(?!main$).*"}, "dev", true},
		{[]string{"^(?!main$).*"}, "main", false},
	}
	for _, c := range cases {
		if got := matchAnyRegex(c.patterns, c.branch); got != c.want {
			t.Errorf("MatchAnyRegex(%v, %q) = %v, want %v", c.patterns, c.branch, got, c.want)
		}
	}
}
