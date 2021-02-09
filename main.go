package main

import (
	"bytes"
	"github.com/spf13/cobra"
	c "go.amplifyedge.org/booty-v2/cmd"
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/pkg/logging/zaplog"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

//golang 1.16 should be able to do this natively
//go:embed config.json


const defaultConfig = `
{
  "dev": true,
  "binaries": [
    {
      "name": "grafana",
      "version": "7.4.0"
    },
    {
      "name": "goreleaser",
      "version": "0.155.1"
    },
    {
      "name": "caddy",
      "version": "2.3.0"
    },
    {
      "name": "protoc",
      "version": "3.14.0"
    },
    {
      "name": "protoc-gen-go",
      "version": "1.25.0"
    },
    {
      "name": "protoc-gen-cobra",
      "version": "0.4.0"
    },
    {
      "name": "protoc-gen-go-grpc",
      "version": "master"
    }
  ]
}
`

var (
	rootCmd = &cobra.Command{
		Use: "booty [commands]",
	}
	logger *zaplog.ZapLogger
)

func main() {
	var vi *config.VersionInfo
	var comps []dep.Component
	logger = zaplog.NewZapLogger(zaplog.WARN, "booty", true)
	logger.InitLogger(nil)
	// setup directories
	err := osutil.SetupDirs()
	if err != nil {
		logger.Fatalf("unable to setup directories: %v", err)
	}
	// config, loads default config if it doesn't exists
	etc := osutil.GetEtcDir()
	var r io.Reader
	fileContent, err := ioutil.ReadFile(filepath.Join(etc, "config.json"))
	if err != nil {
		// use default config
		r = strings.NewReader(defaultConfig)
	} else {
		r = bytes.NewBuffer(fileContent)
	}
	vi = config.NewVersionInfo(logger, r)

	// global db directory
	db := store.NewDB(logger, osutil.GetDataDir())
	comps = append(comps,
		components.NewCaddy(db, vi.GetVersion("caddy")),
	)
	if vi.DevMode {
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
			c.ReleaseCommand(logger, comps),
			c.ProtoCommand(logger, comps),
		)
	}
	rootCmd.AddCommand(c.InstallCommand(logger, comps))
	if err = rootCmd.Execute(); err != nil {
		logger.Errorf("error: %v", err)
	}
}
