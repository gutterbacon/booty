package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
)

// InstallAll installs the component's binary to prefix
func InstallAllCommand(e dep.Executor) *cobra.Command {
	installAllCmd := &cobra.Command{Use: "install-all"}
	installAllCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := e.DownloadAll(); err != nil {
			return err
		}
		return e.InstallAll()
	}
	return installAllCmd
}
