// package depmanager manages dependencies of a running system based on its
// specified environment (i.e. dev or user)
package dep

import (
	"github.com/spf13/cobra"
	"go.amplifyedge.org/booty-v2/internal/gitutil"
	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/update"
)

// Executor is responsible for
type Executor interface {
	Component(name string) Component
	AllComponents() []Component              // get all components for an env (dev or user for example)
	AllInstalledComponents() ([]byte, error) // list all installed components
	DownloadAll() error                      // fetch all components
	Run(name string, args ...string) error
	RunAll() error // run all service components
	Stop(name string) error
	StopAll() error
	Install(name, version string) error // installs a single component by its name
	InstallAll() error                  // install all components
	Uninstall(name string) error        // uninstalls a component
	UninstallAll() error                // uninstall all components
	Backup(name string) error           // backup single component by its name
	BackupAll() error                   // backup all components
	CleanCache() error
}

// Agent is responsible for
type Agent interface {
	Checker() *update.Checker // do work as agent (view updates, collect metrics if any etc), returns status code of the operation
}

// Commander has to be able to output a cobra.Command and logging.Logger
type Commander interface {
	Logger() logging.Logger
	Command() *cobra.Command
	Completion() ([]byte, error)
}

// Extractor extracts embedded files to some determined place
type Extractor interface {
	Extract(string) error // extract makefiles to a directory
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
	Version() update.Version
	SetVersion(update.Version)
	Download() error // download to dir
	Dependencies() []Component
	Install() error
	Uninstall() error
	Run(args ...string) error
	Update(version update.Version) error
	RunStop() error
	Backup() error
	IsDev() bool
	IsService() bool
	RepoUrl() update.RepositoryURL
}

// replacing os.mk and help.mk
type OSPrinter interface {
	OSInfo() string
}

// replacing gitr.mk
type GitWrapper interface {
	SetupFork(upstreamOwner string) error
	CatchupFork() error
	CatchupAll() error
	RegisterRepos(directories ...string) error
	StageAll() error
	Stage(args ...string) error
	Commit(msg string) error
	Push() error
	SubmitPR() error
	CreateTag(tagName string, tagMsg string) error
	PushTag() error
	DeleteTag(tagName string) error
	RepoInfo(dirpath string) (*gitutil.RepoInfo, error)
}
