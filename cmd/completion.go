package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"os"
)

func CompletionCommand(c dep.Commander) *cobra.Command {
	const desc = "generate completion for current shell, outputs it to stdout. Have a look at your shell documentation to install it."
	complCmd := &cobra.Command{
		Use:   "completion",
		Short: desc,
		Long:  desc,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := c.Completion()
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(os.Stdout, string(b))
			return err
		},
	}
	return complCmd
}
