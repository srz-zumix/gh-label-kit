package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/cmd/repo"
)

func NewRepoCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "repo",
		Short: "Manage repository labels",
		Long:  `Manage repository labels.`,
	}

	cmd.AddCommand(repo.NewCopyCmd())
	cmd.AddCommand(repo.NewListCmd())
	cmd.AddCommand(repo.NewSyncCmd())

	return cmd
}

func init() {
	rootCmd.AddCommand(NewRepoCmd())
}
