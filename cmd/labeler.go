package cmd

import (
	"context"
	"fmt"
	"os"

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
	var nameOnly bool
	var syncLabels bool
	var dryrun bool
	var ref string
	cmd := &cobra.Command{
		Use:   "labeler <pr-number...>",
		Short: "Automatically label PRs based on changed files and branch name using config file",
		Long:  `Automatically add or remove labels to GitHub Pull Requests based on changed files, branch name, and a YAML config. Supports glob/regex patterns and syncLabels option for label removal. https://github.com/actions/labeler`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repository, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing repository: %w", err)
			}

			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("error creating GitHub client: %w", err)
			}

			cfg, err := labeler.LoadConfig(configPath)
			if err != nil {
				if repo == "" {
					return fmt.Errorf("failed to load config: %w", err)
				}
				if ref == "" {
					ref = os.Getenv("GITHUB_SHA")
				}
				cfg, err = labeler.LoadConfigFromRepo(ctx, client, repository, configPath, &ref)
				if err != nil {
					return fmt.Errorf("failed to load config from repository: %w", err)
				}
			}

			for _, prNumber := range args {
				pr, err := gh.GetPullRequest(ctx, client, repository, prNumber)
				if err != nil {
					return fmt.Errorf("failed to get PR %s: %w", prNumber, err)
				}

				changedFiles, err := gh.ListPullRequestFiles(ctx, client, repository, prNumber)
				if err != nil {
					return fmt.Errorf("failed to get PR files for %s: %w", prNumber, err)
				}

				result := labeler.CheckMatchConfigs(cfg, changedFiles, pr)
				allLabels := result.GetLabels(syncLabels)

				if dryrun {
					if result.HasDiff(syncLabels) {
						fmt.Printf("Would set labels for PR #%s: %v to %v\n", prNumber, result.Current, allLabels)
					} else {
						fmt.Printf("No label changes for PR #%s: %v\n", prNumber, allLabels)
					}
				} else {
					renderer := render.NewRenderer(opts.Exporter)
					labels := pr.Labels
					if result.HasDiff(syncLabels) {
						renderer.WriteLine(fmt.Sprintf("Labels set for PR #%s", prNumber))
						labels, err = labeler.SetLabels(ctx, client, repository, pr, allLabels)
						if err != nil {
							return fmt.Errorf("failed to set labels for PR %s: %w", prNumber, err)
						}
					} else {
						renderer.WriteLine(fmt.Sprintf("No label changes for PR #%s", prNumber))
					}
					renderer.SetColor(colorFlag)
					if nameOnly {
						renderer.RenderNamesWithSeparator(labels, ",")
					} else {
						renderer.RenderLabelsDefault(labels)
					}
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", render.ColorFlagAuto, render.ColorFlags, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Target repository in the format 'owner/repo'")
	f.StringVar(&configPath, "config", ".github/labeler.yml", "Path to labeler config YAML file")
	f.BoolVar(&nameOnly, "name-only", false, "Output only team names")
	f.BoolVar(&syncLabels, "sync", false, "Remove labels not matching any condition")
	f.BoolVarP(&dryrun, "dryrun", "n", false, "Dry run: do not actually set labels")
	f.StringVar(&ref, "ref", "", "Git reference (branch, tag, or commit SHA) to load config from repository")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}

func init() {
	rootCmd.AddCommand(NewLabelerCmd())
}
