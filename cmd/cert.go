package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func CertCommand(e dep.Executor) *cobra.Command {
	runCmd := &cobra.Command{Use: "cert", DisableFlagParsing: true, Short: "run mkcert under the hood"}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Run("mkcert", args...)
	}
	return runCmd
}
