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
		{[]string{"^(?!main$).*"}, "dev", true},
		{[]string{"^(?!main$).*"}, "main", false},
		{[]string{"^(?!ci/)(?!release/).*"}, "feature/abc", true},
		{[]string{"^(?!ci/)(?!release/).*"}, "ci/abc", false},
		{[]string{"^(?!ci/)(?!release/).*"}, "release/abc", false},
	}
	for _, c := range cases {
		if got := matchAnyRegex(c.patterns, c.branch); got != c.want {
			t.Errorf("MatchAnyRegex(%v, %q) = %v, want %v", c.patterns, c.branch, got, c.want)
		}
	}
}
