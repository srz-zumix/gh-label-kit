package repo

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/go-gh-extension/pkg/gh"
	"github.com/srz-zumix/go-gh-extension/pkg/parser"
)

func NewSyncCmd() *cobra.Command {
	var repo string

	cmd := &cobra.Command{
		Use:   "sync <dst-repository...>",
		Short: "Sync labels from source repository to destination repository",
		Long:  `Sync all labels from the source repository to the destination repositories. If a label already exists in the destination, it will be updated if --force is specified.`,
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
				if err := gh.SyncLabels(ctx, client, src, dst); err != nil {
					return fmt.Errorf("failed to sync labels to %s: %w", dstArg, err)
				}
				fmt.Printf("Successfully synced labels from %s to %s\n", args[0], dstArg)
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repo, "repo", "R", "", "The repository in the format 'owner/repo'")

	return cmd
}
