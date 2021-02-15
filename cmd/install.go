package cmd

import (
	"errors"

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

// Install installs the component with name as args
func InstallCommand(e dep.Executor) *cobra.Command {
	installCmd := &cobra.Command{
		Use:  "install <name> <version>",
		Args: cobra.ExactArgs(2),
	}
	installCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("use install <component_name> <desired_version>")
		}
		return e.Install(args[0], args[1])
	}
	return installCmd
}
