package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"os"
)

func ListAllCommand(e dep.Executor) *cobra.Command {
	installAllCmd := &cobra.Command{Use: "list-installed", Short: "list all installed components"}
	installAllCmd.RunE = func(cmd *cobra.Command, args []string) error {
		b, err := e.AllInstalledComponents()
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(os.Stdout, b)
		return err
	}
	return installAllCmd
}
