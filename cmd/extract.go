package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func ExtractCommand(e dep.Extractor) *cobra.Command {
	extractCmd := &cobra.Command{
		Use:                "extract <DESTINATION>",
		Short:              "extract makefiles to directory of choice",
		Example:            `extract "./here"`,
		Args:               cobra.ExactArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return e.Extract(args[0])
		},
	}
	extractCmd.DisableFlagParsing = true
	extractCmd.Flags().SetInterspersed(true)
	return extractCmd
}
