/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/srz-zumix/gh-label-kit/version"
	"github.com/srz-zumix/go-gh-extension/pkg/actions"
)

var rootCmd = &cobra.Command{
	Use:     "gh-label-kit",
	Short:   "A tool to manage GitHub labels",
	Long:    `gh-label-kit is a tool to manage GitHub labels.`,
	Version: version.Version,
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
}
