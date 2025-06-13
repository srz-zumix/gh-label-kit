package runner

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type RunnerListOptions struct {
	Exporter cmdutil.Exporter
}

func NewRunnerListCmd() *cobra.Command {
	opts := &RunnerListOptions{}
	var owner string
	var repo string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List GitHub Actions runner labels in a repository",
		Long:  `List all GitHub Actions runner labels in the specified repository.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := parser.Repository(parser.RepositoryOwner(owner), parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}
			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(r)
			if err != nil {
				return fmt.Errorf("error creating GitHub client: %w", err)
			}
			runners, err := gh.ListRunners(ctx, client, r)
			if err != nil {
				return fmt.Errorf("failed to list runner labels: %w", err)
			}
			renderer := render.NewRenderer(opts.Exporter)
			renderer.RenderRunnersDefault(runners)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}
