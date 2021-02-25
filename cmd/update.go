package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func UpdateAllCommand(a dep.Agent, c dep.Commander) *cobra.Command {
	updCmd := &cobra.Command{Use: "update-all", Short: "update-all components including booty itself", Args: cobra.NoArgs}
	updCmd.RunE = func(cmd *cobra.Command, args []string) error {
		checker := a.Checker()
		c.Logger().Info("checking all components update")
		if err := checker.CheckNewReleases(); err != nil {
			return err
		}
		return nil
	}
	return updCmd
}
