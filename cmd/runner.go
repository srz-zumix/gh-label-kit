package cmd

import (
	"github.com/spf13/cobra"

	"github.com/srz-zumix/gh-label-kit/cmd/runner"
)

func NewRunnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Manage runner labels",
		Long:  `Manage GitHub Actions runner labels.`,
	}
	cmd.AddCommand(runner.NewRunnerListCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(NewRunnerCmd())
}
