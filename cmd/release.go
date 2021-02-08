package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/pkg/logging"
)

// runs Goreleaser under the hood
func ReleaseCommand(logger logging.Logger, comps []dep.Component) *cobra.Command {
	runCmd := &cobra.Command{Use: "release", DisableFlagParsing: true}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, c := range comps {
			if c.Name() == "goreleaser" {
				logger.Infof("running goreleaser version: %s", c.Version())
				return c.Run(args...)
			}
		}
		logger.Error("goreleaser not found")
		return nil
	}
	return runCmd
}
