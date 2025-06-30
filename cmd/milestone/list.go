package milestone

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type ListOptions struct {
	Exporter cmdutil.Exporter
}

func NewListCmd() *cobra.Command {
	opts := &ListOptions{}
	var colorFlag string
	var repo string
	cmd := &cobra.Command{
		Use:   "list <milestone>",
		Short: "List labels for a milestone",
		Long:  `List all labels attached to issues and PRs in the specified milestone.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			milestone := args[0]
			milestoneNum, err := strconv.Atoi(milestone)
			if err != nil || milestoneNum <= 0 {
				return fmt.Errorf("invalid milestone number: %s", milestone)
			}
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}
			ctx := context.Background()
			labels, err := gh.ListLabelsForMilestone(ctx, client, repository, milestoneNum)
			if err != nil {
				return fmt.Errorf("failed to list labels for milestone '%s': %w", milestone, err)
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
