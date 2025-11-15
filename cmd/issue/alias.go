package issue

import (
	"os"

	"github.com/spf13/cobra"
)

// IsCalledViaAlias checks if the parent command was called via an alias
// by examining the actual command line arguments
func IsCalledViaAlias(cmd *cobra.Command, aliasName string) bool {
	parentCmd := cmd.Parent()
	if parentCmd == nil {
		return false
	}

	args := os.Args
	// Find the parent command name in args
	// The structure should be: [program] [parent-cmd] [sub-cmd] [flags/args...]
	if len(args) > 1 {
		// args[1] should be the parent command (issue or pr)
		parentArg := args[1]
		return parentArg == aliasName
	}

	return false
}

// IsPRCommand checks if the command was called via the "pr" alias
func IsPRCommand(cmd *cobra.Command) bool {
	return IsCalledViaAlias(cmd, "pr")
}
