package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
)

// wrapper around protoc
// runs protoc under the hood with the necessary include
func ProtoCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "protoc", DisableFlagParsing: true, Short: "wrapper for protoc"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("protoc", args...)
	}
	return runCmd
}
