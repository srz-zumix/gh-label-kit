package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/cmd/discussion"
)

// NewDiscussionCmd creates the discussion command and registers subcommands
func NewDiscussionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discussion",
		Short: "Manage discussion labels",
		Long:  `Manage discussion labels.`,
	}

	cmd.AddCommand(discussion.NewListCmd())
	cmd.AddCommand(discussion.NewAddCmd())
	cmd.AddCommand(discussion.NewRemoveCmd())
	cmd.AddCommand(discussion.NewClearCmd())
	cmd.AddCommand(discussion.NewSetCmd())
	cmd.AddCommand(discussion.NewSearchCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewDiscussionCmd())
}
