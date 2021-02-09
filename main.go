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
)

const (
	defaultDev             = false
	defaultVersionInfoFile = "./components_version.json"
)

var (
	isDev       bool
	versionInfo string
)

func main() {
	logger := zaplog.NewZapLogger(zaplog.INFO, "booty", true)
	logger.InitLogger(nil)
	// setup directories
	err := osutil.SetupDirs()
	if err != nil {
		logger.Fatalf("unable to setup directories: %v", err)
	}

	// global db directory
	db := store.NewDB(logger, osutil.GetDataDir())
	rootCmd := &cobra.Command{Use: "booty [commands]"}
	rootCmd.PersistentFlags().BoolVarP(&isDev, "dev", "d", defaultDev, "run tools as dev activating several more components useful for developing")
	rootCmd.PersistentFlags().StringVarP(&versionInfo, "config-version-info-file", "c", defaultVersionInfoFile, "path to config file")
	var vi *config.VersionInfo
	var comps []dep.Component
	vi = config.NewVersionInfo(logger, versionInfo)
	comps = append(comps,
		components.NewCaddy(db, vi.GetVersion("caddy")),
	)
	if isDev {
		if err = osutil.DetectPreq(); err != nil {
			logger.Fatalf(err.Error())
		}
		comps = append(comps,
			components.NewGoreleaser(db, vi.GetVersion("goreleaser")),
			components.NewProtocGenGo(db, vi.GetVersion("protoc-gen-go")),
			//components.NewProtocGenGoGrpc(db, vi.GetVersion("protoc-gen-go-grpc")),
			components.NewProtocGenCobra(db, vi.GetVersion("protoc-gen-cobra")),
			components.NewProtoc(db, vi.GetVersion("protoc")),
			components.NewGrafana(db, vi.GetVersion("grafana")),
		)
		rootCmd.AddCommand(
			cmd.ReleaseCommand(logger, comps),
			cmd.ProtoCommand(logger, comps),
		)
	}
	rootCmd.AddCommand(cmd.InstallCommand(logger, comps))

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("error: %v", err)
	}
}
