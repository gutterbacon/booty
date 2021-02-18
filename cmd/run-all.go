package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

// runs all services under the hood
func RunAllCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "run-all", DisableFlagParsing: true}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.RunAll()
	}
	return runCmd
}
