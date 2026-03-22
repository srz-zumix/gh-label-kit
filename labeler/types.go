package labeler

import "github.com/google/go-github/v84/github"

type PullRequest = github.PullRequest
type CommitFile = github.CommitFile
type Label = github.Label
type User = github.User
type PullRequestBranch = github.PullRequestBranch

func Ptr[T any](v T) *T {
	return &v
}
