package orchestrator

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"

	"go.amplifyedge.org/booty-v2/cmd"
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
)

// Orchestrator implements Agent
type Orchestrator struct {
	cfg        *config.AppConfig
	components map[string]dep.Component
	logger     logging.Logger
	command    *cobra.Command
}

// bootstraps everything
func NewOrchestrator(app string) *Orchestrator {
	var ac *config.AppConfig
	comps := map[string]dep.Component{}
	var rootCmd = &cobra.Command{
		Use: app,
	}
	// setup logger
	logger := zaplog.NewZapLogger(zaplog.WARN, app, true)
	logger.InitLogger(nil)
	// config, loads default config if it doesn't exists
	etc := osutil.GetEtcDir()
	configPath := filepath.Join(etc, "config.json")
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		// use default config
		fileContent, err = json.Marshal(&config.DefaultConfig)
		if err != nil {
			logger.Fatalf("error encoding default json config")
		}
		// write it to default path
		err = ioutil.WriteFile(configPath, fileContent, 0644)
		if err != nil {
			logger.Fatalf("error writing default json config")
		}
	}
	ac = config.NewAppConfig(logger, fileContent)

	// setup badger database for package tracking
	db := store.NewDB(logger, osutil.GetDataDir())
	comps["caddy"] = components.NewCaddy(db, ac.GetVersion("caddy"))
	if ac.DevMode {
		if err = osutil.DetectPreq(); err != nil {
			logger.Fatalf(err.Error())
		}
		comps["goreleaser"] = components.NewGoreleaser(db, ac.GetVersion("goreleaser"))
		comps["grafana"] = components.NewGrafana(db, ac.GetVersion("grafana"))
		comps["protoc-gen-go"] = components.NewProtocGenGo(db, ac.GetVersion("protoc-gen-go"))
		comps["protoc-gen-go-grpc"] = components.NewProtocGenGoGrpc(db, ac.GetVersion("protoc-gen-go-grpc"))
		comps["protoc-gen-cobra"] = components.NewProtocGenCobra(db, ac.GetVersion("protoc-gen-cobra"))
		comps["protoc"] = components.NewProtoc(
			db, ac.GetVersion("protoc"),
			[]dep.Component{
				components.NewProtocGenGo(db, ac.GetVersion("proto-gen-go")),
				components.NewProtocGenGoGrpc(db, ac.GetVersion("protoc-gen-go-grpc")),
				components.NewProtocGenGoGrpc(db, ac.GetVersion("protoc-gen-cobra")),
			},
		)
	}
	return &Orchestrator{
		cfg:        ac,
		components: comps,
		logger:     logger,
		command:    rootCmd,
	}
}

func (o *Orchestrator) Command() *cobra.Command {
	extraCmds := []*cobra.Command{
		cmd.InstallAllCommand(o),
	}
	if o.cfg.DevMode {
		extraCmds = append(
			extraCmds,
			cmd.ProtoCommand(o),
			cmd.ReleaseCommand(o),
		)
	}
	o.command.AddCommand(extraCmds...)
	return o.command
}

func (o *Orchestrator) Logger() logging.Logger {
	return o.logger
}

func (o *Orchestrator) Component(name string) dep.Component {
	for _, comp := range o.components {
		if comp.Name() == name {
			return comp
		}
	}
	return nil
}

func (o *Orchestrator) Components() []dep.Component {
	var comps []dep.Component
	for _, v := range o.components {
		comps = append(comps, v)
	}
	return comps
}

func (o *Orchestrator) DownloadAll() error {
	o.logger.Info("downloading all components")
	for _, c := range o.components {
		if err := c.Download(osutil.GetDownloadDir()); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) Install(name, version string) error {
	// TODO
	return nil
}

func (o *Orchestrator) InstallAll() error {
	o.logger.Info("installing all components")
	for _, c := range o.components {
		if err := c.Install(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) Backup(name string) error {
	// TODO
	return nil
}

func (o *Orchestrator) BackupAll() error {
	// TODO
	return nil
}

func (o *Orchestrator) Run(name string, args ...string) error {
	comp := o.Component(name)
	return comp.Run(args...)
}
