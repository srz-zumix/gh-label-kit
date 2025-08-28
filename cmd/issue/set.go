package issue

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type SetOptions struct {
	Exporter cmdutil.Exporter
}

func NewSetCmd() *cobra.Command {
	opts := &SetOptions{}
	var colorFlag string
	var repo string
	cmd := &cobra.Command{
		Use:   "set <number> <label>...",
		Short: "Set labels for a issue (replace all)",
		Long:  `Set (replace) all labels for a issue in the repository.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			issue := args[0]
			labels := args[1:]
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			result, err := gh.SetIssueLabels(ctx, client, repository, issue, labels)
			if err != nil {
				return fmt.Errorf("failed to set labels for issue #%s: %w", issue, err)
			}
			renderer := render.NewRenderer(opts.Exporter)
			renderer.SetColor(colorFlag)
			renderer.RenderLabelsDefault(result)
			return nil
		},
	}
	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", render.ColorFlagAuto, render.ColorFlags, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)
	return cmd
}
