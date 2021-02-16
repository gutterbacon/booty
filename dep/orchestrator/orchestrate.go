package orchestrator

import (
	"fmt"
	"github.com/spf13/cobra"
	registry2 "go.amplifyedge.org/booty-v2/dep/registry"
	"go.amplifyedge.org/booty-v2/internal/errutil"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.amplifyedge.org/booty-v2/cmd"
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
)

// Orchestrator implements Executor, Agent, and Commander
type Orchestrator struct {
	cfg        *config.AppConfig
	components map[string]dep.Component
	logger     logging.Logger
	command    *cobra.Command
}

// constructor
func NewOrchestrator(app string) *Orchestrator {
	var ac *config.AppConfig
	var comps map[string]dep.Component
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
		fileContent, err = config.DefaultConfig.ReadFile("config.reference.json")
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
	db := store.NewDB(logger, filepath.Join(osutil.GetDataDir(), "packages"))

	// setup registry
	registry, err := registry2.NewRegistry(db, ac)
	if err != nil {
		logger.Fatalf("error creating components: %v", err)
	}
	if ac.DevMode {
		comps = registry.DevComponents
	} else {
		comps = registry.Components
	}

	return &Orchestrator{
		cfg:        ac,
		components: comps,
		logger:     logger,
		command:    rootCmd,
	}
}

// =================================================================
// Commander
// =================================================================
func (o *Orchestrator) Command() *cobra.Command {
	extraCmds := []*cobra.Command{
		cmd.InstallAllCommand(o),
		cmd.InstallCommand(o),
		cmd.UninstallAllCommand(o),
		cmd.UninstallCommand(o),
		cmd.JbCommand(o),
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

// =================================================================
// Executor
// =================================================================

func (o *Orchestrator) Component(name string) dep.Component {
	o.logger.Debugf("querying component of name: %s", name)
	if o.components[name] != nil {
		return o.components[name]
	}
	return nil
}

func (o *Orchestrator) AllComponents() []dep.Component {
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
		tasks = append(tasks, newTask(k.Download, dlErr(k)))
	}
	pool := newTaskPool(tasks)
	pool.runAll()
	for _, t := range pool.tasks {
		if t.err != nil {
			return t.errFunc(t.err)
		}
	}
	return nil
}

func dlErr(c dep.Component) func(err error) error {
	return func(err error) error {
		return errutil.New(errutil.ErrDownloadComponent, fmt.Errorf("name: %s, version: %s, error: %v", c.Name(), c.Version(), err))
	}
}

func (o *Orchestrator) Run(name string, args ...string) error {
	comp := o.Component(name)
	return comp.Run(args...)
}

func (o *Orchestrator) Install(name, version string) error {
	var err error
	c := o.Component(name)
	if c == nil {
		err = errutil.New(errutil.ErrInvalidComponent, fmt.Errorf("name: %s, version: %s", name, version))
		return err
	}
	// try installing it
	if err = c.Install(); err != nil {
		return errutil.New(errutil.ErrInstallComponent, fmt.Errorf("name: %s, version: %s, err: %v", name, version, err))
	}
	return nil
}

func (o *Orchestrator) InstallAll() error {
	// we don't run concurrently here.
	o.logger.Info("installing all components")
	for _, c := range o.components {
		k := c
		o.logger.Info("installing %s, version: %s", k.Name(), k.Version())
		if err := k.Install(); err != nil {
			return errutil.New(errutil.ErrInstallComponent, fmt.Errorf("name: %s, version: %s, err: %v", k.Name(), k.Version(), err))
		}
	}
	return nil
}

func (o *Orchestrator) Uninstall(name string) error {
	var err error
	o.logger.Info("uninstall %s version %s", name)
	for _, c := range o.components {
		k := c
		if k.Name() == name {
			if err = k.Uninstall(); err != nil {
				return errutil.New(errutil.ErrUninstallComponent, fmt.Errorf("name: %s, version: %s, err: %v", k.Name(), k.Version(), err))
			}
		}
	}
	return nil
}

func (o *Orchestrator) UninstallAll() error {
	o.logger.Info("uninstall all components")
	for _, c := range o.components {
		if err := c.Uninstall(); err != nil {
			return errutil.New(errutil.ErrUninstallComponent, fmt.Errorf("name: %s, version: %s, err: %v", c.Name(), c.Version(), err))
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

func (o *Orchestrator) AllInstalledComponents() []dep.Component {
	// TODO
	return nil
}

// ================================================================
// Agent
// ================================================================

func (o *Orchestrator) Serve() error {
	// TODO run as agent in a separate database to check for updates
	// or to collect metrics for that matter
	var err error
	agentDbDir := filepath.Join(osutil.GetDataDir(), "agent-store")
	err = os.MkdirAll(agentDbDir, 0755)
	if err != nil {
		return err
	}
	//agentDb := store.NewDB(o.logger, agentDbDir)
	return nil
}
