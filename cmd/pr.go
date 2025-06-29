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

	cmd.AddCommand(pr.NewPRLabelsCmd())
	cmd.AddCommand(pr.NewPRAddLabelCmd())
	cmd.AddCommand(pr.NewPRRemoveLabelCmd())
	cmd.AddCommand(pr.NewPRClearLabelCmd())
	cmd.AddCommand(pr.NewPRSetLabelCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewPRCmd())
}
