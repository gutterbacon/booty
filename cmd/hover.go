package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func HoverCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "hover", DisableFlagParsing: true, Short: "wrapper for hover"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("hover", args[1:]...)
	}
	return runCmd
}
