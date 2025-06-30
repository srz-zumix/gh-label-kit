package cmd

import (
	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/cmd/milestone"
)

func NewMilestoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "Manage milestones",
		Long:  `Manage milestones.`,
	}

	cmd.AddCommand(milestone.NewListCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewMilestoneCmd())
}
