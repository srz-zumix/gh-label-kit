package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/cmd/pr"
)

// NewPRCmd creates the pr command and registers subcommands
func NewPRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
		Long:  `Manage pull requests.`,
	}

	cmd.AddCommand(pr.NewListCmd())
	cmd.AddCommand(pr.NewAddCmd())
	cmd.AddCommand(pr.NewRemoveCmd())
	cmd.AddCommand(pr.NewClearCmd())
	cmd.AddCommand(pr.NewSetCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewPRCmd())
}
