package orchestrator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/cmd"
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	registry2 "go.amplifyedge.org/booty-v2/dep/registry"
	"go.amplifyedge.org/booty-v2/internal/errutil"
	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/store/file"
	"go.amplifyedge.org/booty-v2/internal/update"
	sharedCmd "go.amplifyedge.org/shared-v2/tool/bs-crypt/cmd"
	langCmd "go.amplifyedge.org/shared-v2/tool/bs-lang/cmd"
)

const (
	gracefulPeriod = 5 * time.Second
)

// Orchestrator implements Executor, Agent, and Commander
type Orchestrator struct {
	cfg        *config.AppConfig
	components map[string]dep.Component
	logger     logging.Logger
	command    *cobra.Command
	db         store.Storer
}

// constructor
func NewOrchestrator(app string) *Orchestrator {
	var ac *config.AppConfig
	var comps map[string]dep.Component
	var rootCmd = &cobra.Command{
		Use: app,
	}
	// setup logger
	logger := zaplog.NewZapLogger(zaplog.INFO, app, true)
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

	// setup file database for package tracking
	db, err := file.NewDB(logger, filepath.Join(osutil.GetDataDir(), "packages"))
	if err != nil {
		logger.Fatalf("error creating database: %v", err)
	}

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
		db:         db,
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
		cmd.AgentCommand(o),
		cmd.RunAllCommand(o),
		cmd.ListAllCommand(o),
		// here we exported all the internal tools we might need (bs-crypt, bs-lang, etc)
		sharedCmd.EncryptCmd(),
		sharedCmd.DecryptCmd(),
		langCmd.RootCmd,
	}
	if o.cfg.DevMode {
		extraCmds = append(
			extraCmds,
			cmd.ProtoCommand(o),
			cmd.ReleaseCommand(o),
			cmd.JbCommand(o),
			cmd.JsonnetCommand(o),
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
	err := o.setupVersions()
	if err != nil {
		return err
	}
	//var tasks []*task
	//for _, c := range o.components {
	//	k := c
	//	tasks = append(tasks, newTask(k.Download, dlErr(k)))
	//}
	//pool := newTaskPool(tasks)
	//pool.runAll()
	//for _, t := range pool.tasks {
	//	if t.err != nil {
	//		return t.errFunc(t.err)
	//	}
	//}
	for _, c := range o.components {
		if c.Dependencies() != nil {
			for _, d := range c.Dependencies() {
				if err = d.Download(); err != nil {
					return err
				}
			}
		}
		if err = c.Download(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) setupVersions() error {
	var tasks []*task
	for _, c := range o.components {
		k := c
		if k.Dependencies() != nil {
			for _, d := range k.Dependencies() {
				tasks = append(tasks, newTask(setVersion(o.cfg, d), dlErr(d)))
			}
		}
		tasks = append(tasks, newTask(setVersion(o.cfg, k), dlErr(k)))
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
	// download it
	c.SetVersion(update.Version(version))
	if c.Dependencies() != nil {
		for _, d := range c.Dependencies() {
			if err = setVersion(o.cfg, d)(); err != nil {
				return err
			}
			if err = d.Download(); err != nil {
				return err
			}
			if err = d.Install(); err != nil {
				return err
			}
		}
	}
	if err = c.Download(); err != nil {
		return errutil.New(errutil.ErrDownloadComponent, fmt.Errorf("name: %s, version: %s, err: %v", c.Name(), c.Version(), err))
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
		o.logger.Infof("installing %s, version: %s", k.Name(), k.Version())
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
			o.logger.Infof("uninstalling %s, version: %s", k.Name(), k.Version())
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
			o.logger.Infof("uninstalling %s, version: %s", c.Name(), c.Version())
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

func (o *Orchestrator) AllInstalledComponents() ([]byte, error) {
	// TODO
	pkgs, err := o.db.List()
	if err != nil {
		return nil, err
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	tw := tabwriter.NewWriter(buf, 20, 2, 4, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintf(tw, "%s", headerColor("Installed Packages:"))
	for _, p := range pkgs {
		printRow(tw, "\n%s\t\t%s", p.Name, p.Version)
	}
	err = tw.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var (
	headerColor  = color.New(color.FgCyan).SprintFunc()
	pkgColor     = color.New(color.FgHiGreen).SprintFunc()
	versionColor = color.New(color.FgYellow).SprintFunc()
)

func printRow(out io.Writer, format, key string, value string) {
	_, _ = fmt.Fprintf(out, format, pkgColor(key), versionColor(value))
}

func (o *Orchestrator) RunAll() error {
	for _, c := range o.components {
		if c.IsService() {
			if err := c.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ================================================================
// Agent
// ================================================================

func (o *Orchestrator) Serve() int {
	var err error
	// checks twice in one day
	dur := 12 * time.Hour
	// create new checker instance
	repos := map[update.RepositoryURL]update.Version{}
	for _, c := range o.components {
		repos[c.RepoUrl()] = c.Version()
	}
	checker := update.NewChecker(o.logger, repos, func(r update.RepositoryURL, v update.Version) error {
		// find component matching the repo url
		for _, c := range o.components {
			if c.RepoUrl() == r && v != c.Version() {
				if c.IsService() {
					if err = c.RunStop(); err != nil {
						return errutil.New(errutil.ErrUpdateComponent, fmt.Errorf("stopping component service error, %s version %s error: %v", c.Name(), c.Version(), err))
					}
				}
				if err = c.Update(v); err != nil {
					return errutil.New(errutil.ErrUpdateComponent, fmt.Errorf("%s version %s error: %v", c.Name(), c.Version(), err))
				}
			}
		}
		return nil
	})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(dur)
	for {
		select {
		case <-ticker.C:
			o.logger.Info("checking all components update")
			if err = checker.CheckNewReleases(); err != nil {
				o.logger.Error(err)
			}
		case s := <-sigCh:
			o.logger.Warningf("getting signal: %s, terminating gracefully", s.String())
			// TODO do a proper shutdown instead of sleep like this
			return 0
		}
	}
}

func setVersion(ac *config.AppConfig, c dep.Component) func() error {
	var err error
	return func() error {
		v := ac.GetVersion(c.Name())
		if v == "" {
			v, err = update.GetLatestVersion(c.RepoUrl())
			if err != nil {
				return err
			}
		}
		c.SetVersion(update.Version(v))
		return nil
	}
}
