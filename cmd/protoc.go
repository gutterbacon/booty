package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/pkg/logging"
)

// runs protoc under the hood
func ProtoCommand(logger logging.Logger, comps []dep.Component) *cobra.Command {
	runCmd := &cobra.Command{Use: "protoc", DisableFlagParsing: true}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, c := range comps {
			if c.Name() == "protoc" {
				logger.Infof("running protoc version: %s", c.Version())
				return c.Run(args...)
			}
		}
		logger.Error("protoc not found")
		return nil
	}
	return runCmd
}
