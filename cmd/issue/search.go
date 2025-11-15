package issue

import (
	"context"
	"fmt"
	"strings"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type SearchOptions struct {
	Exporter cmdutil.Exporter
}

func NewSearchCmd() *cobra.Command {
	opts := &SearchOptions{}
	var colorFlag string
	var repo string
	var owner string
	var labels []string
	cmd := &cobra.Command{
		Use:   "search [query...]",
		Short: "Search issues by query",
		Long:  `Search issues in the repository using a search query. The query can include label filters and other search criteria.`,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			if len(labels) > 0 {
				for _, label := range labels {
					query += fmt.Sprintf(" label:\"%s\"", label)
				}
			}

			if owner != "" {
				query += " org:" + owner
			}
			if repo != "" {
				query += " repo:" + repo
			}

			// Check if command was called via alias
			if IsPRCommand(cmd) {
				query += " is:pr"
			}

			repository, err := parser.Repository(parser.RepositoryOwner(owner), parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("failed to resolve repository: %w", err)
			}

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("failed to create GitHub client: %w", err)
			}

			ctx := context.Background()
			issues, err := gh.SearchIssues(ctx, client, repository, query)
			if err != nil {
				return fmt.Errorf("failed to search issues: %w", err)
			}

			renderer := render.NewRenderer(opts.Exporter)
			renderer.SetColor(colorFlag)
			renderer.RenderIssuesDefault(issues)
			return nil
		},
	}
	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", render.ColorFlagAuto, render.ColorFlags, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	f.StringVar(&owner, "owner", "", "Specify the organization name")
	f.StringSliceVarP(&labels, "label", "l", []string{}, "Filter issues by labels")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)
	return cmd
}
