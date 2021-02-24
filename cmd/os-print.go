package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"os"
)

func OsPrintCommand(op dep.OSPrinter) *cobra.Command {
	osPrintCmd := &cobra.Command{Use: "os-print", Short: "os-print", Args: cobra.NoArgs}
	osPrintCmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprint(os.Stdout, op.OSInfo())
		return err
	}
	return osPrintCmd
}
