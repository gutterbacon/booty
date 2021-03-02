package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

// runs all services under the hood
func RunAllCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "run-all", DisableFlagParsing: true, Short: "run-all components that can be run as service."}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.RunAll()
	}
	return runCmd
}

func RunCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "run <name_of_component>", DisableFlagParsing: true, Short: "run specified component"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run(args[0], args[1:]...)
	}
	return runCmd
}
