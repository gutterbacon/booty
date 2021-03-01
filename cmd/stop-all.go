package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

// stops all services under the hood
func StopAllCommand(e dep.Executor) *cobra.Command {
	stopCmd := &cobra.Command{Use: "stop-all", DisableFlagParsing: true, Short: "stop-all components that can be run as service."}
	stopCmd.DisableFlagParsing = true
	stopCmd.Flags().SetInterspersed(true)
	stopCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.StopAll()
	}
	return stopCmd
}

// stops service under the hood
func StopCommand(e dep.Executor) *cobra.Command {
	stopCmd := &cobra.Command{Use: "stop <service_name>", DisableFlagParsing: true, Short: "stop components.", Args: cobra.ExactArgs(1)}
	stopCmd.DisableFlagParsing = true
	stopCmd.Flags().SetInterspersed(true)
	stopCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Stop(args[0])
	}
	return stopCmd
}
