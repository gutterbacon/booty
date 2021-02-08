// package depmanager manages dependencies of a running system based on its
// specified environment (i.e. dev or user)
package dep

// Agent is responsible for
type Agent interface {
	Components(env string) []Component  // get all components for an env (dev or user for example)
	DownloadAll() error                 // fetch all components
	Install(name, version string) error // install single component by its name
	InstallAll() error                  // install all components
	Backup(name string) error           // backup single component by its name
	BackupAll() error                   // backup all components
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
	Download(targetDir string) error // download to dir
	Install() error
	Uninstall() error
	Run(args ...string) error
	Update(version string) error
	Stop() error
	Backup() error
}
