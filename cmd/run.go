package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/pkg/logging"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
)

// Run wraps the component's binary inside booty
// useful because we don't want to pollute user's global PATH
// it all should be self-contained.
func RunCommand(logger logging.Logger, comps []dep.Component) *cobra.Command {
	runCmd := &cobra.Command{Use: "run"}
	runCmd.Flags().SetInterspersed(false)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		command := args[0]
		for _, c := range comps {
			if c.Name() == command {
				logger.Infof("running %s version %s", c.Name(), c.Version())
				return osutil.Exec(command, args[1:]...)
			}
		}
		logger.Errorf("no binary with name: %s found", command)
		return nil
	}
	return runCmd
}
