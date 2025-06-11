package cmd

import (
	"context"
	"fmt"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/labeler"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

type LabelerOptions struct {
	Exporter cmdutil.Exporter
}

// NewLabelerCmd implements a command for GitHub PR auto-labeling based on config and PR info.
func NewLabelerCmd() *cobra.Command {
	opts := &LabelerOptions{}
	var colorFlag string
	var repo string
	var configPath string
	var syncLabels bool
	var dryrun bool
	cmd := &cobra.Command{
		Use:   "labeler <pr-number>",
		Short: "Automatically label PRs based on changed files and branch name using config file",
		Long:  `Automatically add or remove labels to GitHub Pull Requests based on changed files, branch name, and a YAML config. Supports glob/regex patterns and syncLabels option for label removal.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := labeler.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}
			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("error creating GitHub client: %w", err)
			}

			prNumber := args[0]
			pr, err := gh.GetPullRequest(ctx, client, repository, prNumber)
			if err != nil {
				return fmt.Errorf("failed to get PR: %w", err)
			}

			// 変更ファイル一覧取得
			changedFiles, err := gh.ListPullRequestFiles(ctx, client, repository, prNumber)
			if err != nil {
				return fmt.Errorf("failed to get PR files: %w", err)
			}

			result := labeler.CheckMatchConfigs(cfg, changedFiles, pr)

			if dryrun {
				fmt.Printf("Would set labels: %v\n", result.GetLabels(syncLabels))
			} else {
				labels, err := labeler.SetLabels(ctx, client, repository, pr, result, syncLabels)
				if err != nil {
					return fmt.Errorf("failed to set labels: %w", err)
				}
				renderer := render.NewRenderer(opts.Exporter)
				renderer.SetColor(colorFlag)
				renderer.RenderLabelsDefault(labels)
			}
			return nil
		},
	}

	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", "auto", []string{"always", "never", "auto"}, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Target repository in the format 'owner/repo'")
	f.StringVar(&configPath, "config", ".github/labeler.yml", "Path to labeler config YAML file")
	f.BoolVar(&syncLabels, "sync", false, "Remove labels not matching any condition")
	f.BoolVarP(&dryrun, "dryrun", "n", false, "Dry run: do not actually set labels")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}

func init() {
	rootCmd.AddCommand(NewLabelerCmd())
}
