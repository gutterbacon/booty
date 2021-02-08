package main

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/cmd"
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/pkg/logging/zaplog"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"
	"os"
	"path/filepath"
)

const (
	defaultDev             = true
	defaultVersionInfoFile = "./components_version.json"
)

var (
	isDev       bool
	versionInfo string
)

func main() {
	logger := zaplog.NewZapLogger(zaplog.INFO, "booty-v2", true)
	logger.InitLogger(nil)
	// global prefix
	prefix := osutil.GetInstallPrefix()
	prefix = filepath.Join(prefix, "booty")
	err := os.MkdirAll(prefix, 0755)
	if err != nil {
		logger.Fatalf("unable to create prefix directory: %v", err)
	}
	// global db directory
	dbDir := filepath.Join(prefix, "data")
	err = os.MkdirAll(dbDir, 0755)
	if err != nil {
		logger.Fatalf("unable to create data directory: %v", err)
	}
	db := store.NewDB(logger, dbDir)
	rootCmd := &cobra.Command{Use: "booty-v2"}
	rootCmd.PersistentFlags().BoolVarP(&isDev, "dev", "d", defaultDev, "run tools in dev mode instead of user mode")
	rootCmd.PersistentFlags().StringVarP(&versionInfo, "config-version-info-file", "c", defaultVersionInfoFile, "path to config file")
	var vi *config.VersionInfo
	var comps []dep.Component
	vi = config.NewVersionInfo(logger, versionInfo)
	comps = []dep.Component{
		components.NewGrafana(db, vi.GetVersion("grafana")),
		components.NewCaddy(db, vi.GetVersion("caddy")),
	}
	if isDev {
		comps = append(comps, components.NewGoreleaser(db, vi.GetVersion("goreleaser")))
	}
	rootCmd.AddCommand(cmd.ComponentsCommand(logger, comps))

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("error: %v", err)
	}
}
