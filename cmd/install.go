package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
)

// InstallAll installs the component's binary to prefix
func InstallAllCommand(a dep.Agent) *cobra.Command {
	installAllCmd := &cobra.Command{Use: "install-all"}
	installAllCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := a.DownloadAll(); err != nil {
			return err
		}
		return a.InstallAll()
	}
	return installAllCmd
}
