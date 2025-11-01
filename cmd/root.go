/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/version"
	"github.com/srz-zumix/go-gh-extension/pkg/actions"
	"github.com/srz-zumix/go-gh-extension/pkg/logger"
)

var (
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:     "gh-label-kit",
	Short:   "A tool to manage GitHub labels",
	Long:    `gh-label-kit is a tool to manage GitHub labels.`,
	Version: version.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.SetLogLevel(logLevel)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if actions.IsRunsOn() {
		rootCmd.SetErrPrefix(actions.GetErrorPrefix())
	}
	logger.AddCmdFlag(rootCmd, rootCmd.PersistentFlags(), &logLevel, "log-level", "L")
}
