// package depmanager manages dependencies of a running system based on its
// specified environment (i.e. dev or user)
package dep

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/internal/logging"
)

// Executor is responsible for
type Executor interface {
	Component(name string) Component
	AllComponents() []Component // get all components for an env (dev or user for example)
	DownloadAll() error         // fetch all components
	Run(name string, args ...string) error
	Install(name, version string) error // installs a single component by its name
	InstallAll() error                  // install all components
	Uninstall(name string) error        // uninstalls a component
	UninstallAll() error                // uninstall all components
	Backup(name string) error           // backup single component by its name
	BackupAll() error                   // backup all components
}

// Agent is responsible for
type Agent interface {
	Serve() error // do work as agent (view updates, collect metrics if any etc)
}

// Commander has to be able to output a cobra.Command and logging.Logger
type Commander interface {
	Logger() logging.Logger
	Command() *cobra.Command
}

// Component is an interface
// each component has to be able to
// a. Download (along with its 3rd party dependencies)
// b. Install to somewhere that is discoverable by the OS's PATH
// c. Uninstall cleanly
// d. Run itself
// e. Backup its configuration
// has to be able to run each operations on each targeted OSes (Windows, Linux, Mac)
// and each targeted architecture (ARM64, x86_64)
// doesn't matter the implementation whether it will be just makefiles or shell script
// this way if you want to add another component you only have
type Component interface {
	Name() string
	Version() string
	Download() error // download to dir
	Dependencies() []Component
	Install() error
	Uninstall() error
	Run(args ...string) error
	Update(version string) error
	RunStop() error
	Backup() error
}
