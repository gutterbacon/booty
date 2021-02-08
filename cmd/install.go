package cmd

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/pkg/logging"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
)

// Install installs the component's binary to prefix
func InstallCommand(logger logging.Logger, comps []dep.Component) *cobra.Command {
	installCmd := &cobra.Command{Use: "install"}
	installCmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, component := range comps {
			logger.Infof("downloading %s version %s", component.Name(), component.Version())
			if err := component.Download(osutil.GetDownloadDir()); err != nil {
				logger.Debugf("found error while downloading %s component: %v", component.Name(), err)
				return err
			}
			if err := component.Install(); err != nil {
				logger.Debugf("found error while installing %s component: %v", component.Name(), err)
				return err
			}
		}
		return nil
	}
	return installCmd
}
