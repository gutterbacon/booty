package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

// UninstallAll uninstalls all of the components
func UninstallAllCommand(e dep.Executor) *cobra.Command {
	uninstallAllCmd := &cobra.Command{Use: "uninstall-all",  Short: "uninstall all components"}
	uninstallAllCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.UninstallAll()
	}
	return uninstallAllCmd
}

func UninstallCommand(e dep.Executor) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use: "uninstall <name>",
		Short: "uninstall <component_name>",
		Args: cobra.ExactArgs(1),
	}
	uninstallCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return e.Uninstall(args[0])
	}
	return uninstallCmd
}
