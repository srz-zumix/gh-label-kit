package issue

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func NewClearCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "clear <number>",
		Short: "Remove all labels from a issue",
		Long:  `Remove all labels from a issue in the repository.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			repository, err := parser.Repository(parser.RepositoryInput(repo), parser.RepositoryFromURL(target))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			err = gh.ClearIssueLabels(ctx, client, repository, target)
			if err != nil {
				return fmt.Errorf("failed to clear labels from issue %s: %w", target, err)
			}
			logger.Info("All labels removed from issue", "issue", target, "repository", parser.GetRepositoryFullName(repository))
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	return cmd
}
