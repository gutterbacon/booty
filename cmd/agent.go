package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
)

func AgentCommand(a dep.Agent) *cobra.Command {
	agentCmd := &cobra.Command{Use: "agent", Short: "run as agent"}
	agentCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if exitCode := a.Serve(); exitCode != 0 {
			return fmt.Errorf("exit code: %d", exitCode)
		}
		return nil
	}
	return agentCmd
}
