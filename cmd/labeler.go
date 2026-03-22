package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/labeler"
	"github.com/srz-zumix/go-gh-extension/pkg/actions"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
	"github.com/srz-zumix/go-gh-extension/pkg/render"
)

var defaultConfigPath = ".github/labeler.yml"

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
	var reviewRequest string
	var skipLocalConfig bool
	var strictConfig bool
	var noHidden bool
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

			// Set no-hidden option for glob matching
			labeler.SetNoHidden(noHidden)

			client, err := gh.NewGitHubClientWithRepo(repository)
			if err != nil {
				return fmt.Errorf("error creating GitHub client: %w", err)
			}

			var cfg labeler.LabelerConfig
			if !skipLocalConfig {
				if labeler.ConfigFileExists(configPath) {
					cfg, err = labeler.LoadConfig(configPath, strictConfig)
					if err != nil {
						// If local config exists but failed to load, return error immediately
						// Don't fallback to remote config in this case
						return fmt.Errorf("failed to load local config: %w", err)
					}
				}
			}

			ctx := cmd.Context()
			// Load config from repository if local config is skipped or doesn't exist
			if cfg == nil {
				if ref == "" {
					ref = os.Getenv("GITHUB_SHA")
				}
				contentPaths, err := parser.ParseContentPath(configPath)
				if err != nil {
					return fmt.Errorf("failed to parse config path: %w", err)
				}
				if contentPaths.Ref == nil {
					if contentPaths.Repo == nil {
						contentPaths.Ref = &ref
						// If ref is still empty and only one PR is specified, use PR's head branch
						if ref == "" && len(args) == 1 {
							pr, err := gh.GetPullRequest(ctx, client, repository, args[0])
							if err != nil {
								return fmt.Errorf("failed to get PR %s to resolve ref: %w", args[0], err)
							}
							if pr != nil && pr.GetHead() != nil && pr.GetHead().GetRef() != "" {
								ref = pr.GetHead().GetRef()
								logger.Info("Using PR head branch as ref", "pr", args[0], "ref", ref)
							}
						}
					}
				}
				if contentPaths.Repo == nil {
					contentPaths.Repo = &repository
				}
				if contentPaths.Path == nil {
					contentPaths.Path = &defaultConfigPath
				}

				cfg, err = labeler.LoadConfigFromRepo(ctx, client, *contentPaths.Repo, *contentPaths.Path, contentPaths.Ref, strictConfig)
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

				matcher := labeler.NewMatcher(ctx, client)
				result := matcher.CheckMatchConfigs(cfg, changedFiles, pr)
				allLabels := result.GetLabels(syncLabels)
				labeledCodeOwners := labeler.NewLabeledCodeOwners(ctx, client, repository, pr, cfg, reviewRequest)
				reviewRequestLabels := labeler.GetReviewRequestTargetLabels(pr, result, reviewRequest, syncLabels)

				if dryrun {
					if result.HasDiff(syncLabels) {
						logger.Info("Would set labels for PR", "pr", prNumber, "current", result.Current, "new", allLabels)
					} else {
						logger.Info("No label changes for PR", "pr", prNumber, "labels", allLabels)
					}
					codeowners := labeledCodeOwners.GetReviewers(reviewRequestLabels)
					if len(codeowners) > 0 {
						logger.Info("Would request reviewers for PR", "pr", prNumber, "reviewers", codeowners)
					}
				} else {
					renderer := render.NewRenderer(opts.Exporter)
					labels := pr.Labels
					if result.HasDiff(syncLabels) {
						renderer.WriteLine(fmt.Sprintf("Labels set for PR #%s", prNumber))
						labels, err = labeler.SetLabels(ctx, client, repository, pr, allLabels, cfg)
						if err != nil {
							return fmt.Errorf("failed to set labels for PR %s: %w", prNumber, err)
						}
					} else {
						_, err = labeler.EditLabelsByConfig(ctx, client, repository, labels, cfg)
						if err != nil {
							return fmt.Errorf("failed to edit labels for PR %s: %w", prNumber, err)
						}
						renderer.WriteLine(fmt.Sprintf("No label changes for PR #%s", prNumber))
					}
					addedReviewers, _, err := labeledCodeOwners.SetReviewers(reviewRequestLabels)
					if err != nil {
						return fmt.Errorf("failed to set reviewers for PR %s: %w", prNumber, err)
					}
					if len(addedReviewers) > 0 {
						renderer.WriteLine(fmt.Sprintf("Requested reviewers for PR #%s: %v", prNumber, addedReviewers))
					}
					renderer.SetColor(colorFlag)
					if nameOnly {
						err = renderer.RenderNamesWithSeparator(labels, ",")
						if err != nil {
							logger.Warn("Failed to render names, falling back to labels", "pr", prNumber, "error", err)
						}
					} else {
						err = renderer.RenderLabels(labels, nil)
						if err != nil {
							logger.Warn("Failed to render labels, falling back to names", "pr", prNumber, "error", err)
						}
					}
				}

				err = actions.Output("new-labels", strings.Join(result.AddTo(), ","))
				if err != nil {
					return fmt.Errorf("failed to set action output: %w", err)
				}
				err = actions.Output("all-labels", strings.Join(allLabels, ","))
				if err != nil {
					return fmt.Errorf("failed to set action output: %w", err)
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	cmdutil.StringEnumFlag(cmd, &colorFlag, "color", "", render.ColorFlagAuto, render.ColorFlags, "Use color in diff output")
	f.StringVarP(&repo, "repo", "R", "", "Target repository in the format 'owner/repo'")
	f.StringVar(&configPath, "config", defaultConfigPath, "Path to labeler config YAML file, path in repo, or GitHub URL, or actions format (owner/repo[/path]@ref)")
	f.BoolVar(&nameOnly, "name-only", false, "Output only team names")
	f.BoolVar(&syncLabels, "sync", false, "Remove labels not matching any condition")
	f.BoolVarP(&dryrun, "dryrun", "n", false, "Dry run: do not actually set labels")
	f.StringVar(&ref, "ref", "", "Git reference (branch, tag, or commit SHA) to load config from repository")
	f.BoolVar(&skipLocalConfig, "skip-local-config", false, "Skip loading config from local file and load from repository instead")
	f.BoolVar(&strictConfig, "strict", false, "Treat unknown fields in config as errors instead of warnings")
	f.BoolVar(&noHidden, "no-hidden", false, "Exclude hidden files (files starting with .) from glob matching")
	cmdutil.StringEnumFlag(cmd, &reviewRequest, "review-request", "", labeler.ReviewRequestModeAddTo, labeler.ReviewersRequestModes, "Control review request behavior based on CODEOWNERS when labels are applied")
	cmdutil.AddFormatFlags(cmd, &opts.Exporter)

	return cmd
}

func init() {
	rootCmd.AddCommand(NewLabelerCmd())
}
