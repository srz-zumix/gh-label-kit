package pr

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func NewClearCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "clear <pr-number>",
		Short: "Remove all labels from a pull request",
		Long:  `Remove all labels from a pull request in the repository.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pullRequest := args[0]
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			err = gh.ClearPullRequestLabels(ctx, client, repository, pullRequest)
			if err != nil {
				return fmt.Errorf("failed to clear labels from pull request #%s: %w", pullRequest, err)
			}
			fmt.Printf("All labels removed from pull request #%s in repository %s/%s\n", pullRequest, repository.Owner, repository.Name)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	return cmd
}
