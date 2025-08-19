package repo

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func NewCopyCmd() *cobra.Command {
	var force bool
	var repo string

	cmd := &cobra.Command{
		Use:   "copy <dst-repository...>",
		Short: "Copy labels from source repository to destination repository",
		Long:  `Copy all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be skipped unless --force is specified.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := parser.Repository(parser.RepositoryInput(repo))
			if err != nil {
				return fmt.Errorf("error parsing source repository: %w", err)
			}
			ctx := context.Background()
			client, err := gh.NewGitHubClientWithRepo(src)
			if err != nil {
				return fmt.Errorf("error creating GitHub client: %w", err)
			}

			for _, dstArg := range args {
				dst, err := parser.Repository(parser.RepositoryInput(dstArg))
				if err != nil {
					return fmt.Errorf("error parsing destination repository: %w", err)
				}
				if src.Host != dst.Host {
					return fmt.Errorf("source and destination repositories must be on the same host: %s vs %s", src.Host, dst.Host)
				}

				if err := gh.CopyLabels(ctx, client, src, dst, force); err != nil {
					return fmt.Errorf("failed to copy labels to %s: %w", parser.GetRepositoryFullName(dst), err)
				}
				fmt.Printf("Successfully copied labels from %s to %s\n", parser.GetRepositoryFullName(src), parser.GetRepositoryFullName(dst))
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "Repository in the format 'owner/repo'")
	f.BoolVarP(&force, "force", "f", false, "Overwrite existing labels in the destination repository")

	return cmd
}
