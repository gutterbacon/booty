package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func AgentCommand(a dep.Agent, c dep.Commander) *cobra.Command {
	agentCmd := &cobra.Command{Use: "agent", Short: "run as agent"}
	agentCmd.RunE = func(cmd *cobra.Command, args []string) error {
		checker := a.Checker()
		// checks twice in one day
		dur := 12 * time.Hour
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(dur)
		for {
			select {
			case <-ticker.C:
				c.Logger().Info("checking all components update")
				if err := checker.CheckNewReleases(); err != nil {
					c.Logger().Error(err)
				}
			case s := <-sigCh:
				c.Logger().Warningf("getting signal: %s, terminating gracefully", s.String())
				// TODO do a proper shutdown instead of sleep like this
				return nil
			}
		}
	}
	return agentCmd
}
