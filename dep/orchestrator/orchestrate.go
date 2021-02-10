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
	logger := zaplog.NewZapLogger(zaplog.DEBUG, app, true)
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
	cdyVer := ac.GetVersion("caddy")
	gorVer := ac.GetVersion("goreleaser")
	grafVer := ac.GetVersion("grafana")
	protoGoVer := ac.GetVersion("protoc-gen-go")
	protoGrpcVer := ac.GetVersion("protoc-gen-go-grpc")
	protoCobraVer := ac.GetVersion("protoc-gen-cobra")

	// setup badger database for package tracking
	db := store.NewDB(logger, osutil.GetDataDir())
	comps["caddy"] = components.NewCaddy(db, cdyVer)
	if ac.DevMode {
		if err = osutil.DetectPreq(); err != nil {
			logger.Fatalf(err.Error())
		}
		comps["goreleaser"] = components.NewGoreleaser(db, gorVer)
		comps["grafana"] = components.NewGrafana(db, grafVer)
		protoGenGo := components.NewProtocGenGo(db, protoGoVer)
		protoGenGrpc := components.NewProtocGenGoGrpc(db, protoGrpcVer)
		protoCobra := components.NewProtocGenCobra(db, protoCobraVer)
		comps["protoc"] = components.NewProtoc(
			db, ac.GetVersion("protoc"),
			[]dep.Component{
				protoGenGo,
				protoGenGrpc,
				protoCobra,
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
	var tasks []*task
	for _, c := range o.components {
		k := c
		tasks = append(tasks, newTask(k.Download))
	}
	pool := newTaskPool(tasks)
	pool.runAll()
	for _, t := range pool.tasks {
		if t.err != nil {
			return t.err
		}
	}
	return nil
}

func (o *Orchestrator) Install(name, version string) error {
	// TODO
	return nil
}

func (o *Orchestrator) InstallAll() error {
	// we don't run concurrently here.
	o.logger.Info("installing all components")
	for _, c := range o.components {
		k := c
		if err := k.Install(); err != nil {
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

func (o *Orchestrator) UninstallAll() error {
	o.logger.Info("uninstall all components")
	for _, c := range o.components {
		if err := c.Uninstall(); err != nil {
			return err
		}
	}
	return nil
}
