package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func HugoCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "hugo", DisableFlagParsing: true, Short: "wrapper for hugo"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("hugo", args[1:]...)
	}
	return runCmd
}
