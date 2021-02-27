package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func CleanCacheCmd(o dep.Executor) *cobra.Command {
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "clean",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.CleanCache()
		},
	}
	return cleanCmd
}
