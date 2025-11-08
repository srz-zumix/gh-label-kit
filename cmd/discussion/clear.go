package discussion

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
		Short: "Remove all labels from a discussion",
		Long:  `Remove all labels from a discussion in the repository.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			discussion := args[0]
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			err = gh.ClearDiscussionLabels(ctx, client, repository, discussion)
			if err != nil {
				return fmt.Errorf("failed to clear labels from discussion #%s: %w", discussion, err)
			}
			logger.Info("All labels removed from discussion", "discussion", discussion, "repository", parser.GetRepositoryFullName(repository))
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	return cmd
}
