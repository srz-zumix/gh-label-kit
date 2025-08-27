package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/cmd/issue"
)

// NewIssueCmd creates the issue command and registers subcommands
func NewIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Aliases: []string{"pr"},
		Short:   "Manage issue",
		Long:    `Manage issue.`,
	}

	cmd.AddCommand(issue.NewListCmd())
	cmd.AddCommand(issue.NewAddCmd())
	cmd.AddCommand(issue.NewRemoveCmd())
	cmd.AddCommand(issue.NewClearCmd())
	cmd.AddCommand(issue.NewSetCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewIssueCmd())
}
