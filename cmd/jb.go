// wrapper around jb
package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func JbCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "jb", DisableFlagParsing: true, Short: "jsonnet-bundler"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("jb", args...)
	}
	return runCmd
}

func JsonnetCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "jsonnet", DisableFlagParsing: true, Short: "jsonnet"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("jsonnet", args...)
	}
	return runCmd
}
