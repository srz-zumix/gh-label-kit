package discussion

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type RemoveOptions struct {
	Exporter cmdutil.Exporter
}

func NewRemoveCmd() *cobra.Command {
	opts := &RemoveOptions{}
	var colorFlag string
	var repo string
	cmd := &cobra.Command{
		Use:   "remove <number> <label>...",
		Short: "Remove label(s) from a discussion",
		Long:  `Remove one or more labels from a discussion in the repository.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			discussion := args[0]
			labelsToRemove := args[1:]
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			labels, err := gh.RemoveDiscussionLabels(ctx, client, repository, discussion, labelsToRemove)
			if err != nil {
				return fmt.Errorf("failed to remove labels from discussion #%s: %w", discussion, err)
			}
			renderer := render.NewRenderer(opts.Exporter)
			renderer.SetColor(colorFlag)
			renderer.RenderLabelsDefault(labels)
			return nil
		},
	}
	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", render.ColorFlagAuto, render.ColorFlags, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)
	return cmd
}
