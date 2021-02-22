package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	// "path/filepath"

	"go.amplifyedge.org/booty-v2/internal/osutil"
)

const (
	goreleaserBaseRepo = "https://github.com/goreleaser/goreleaser"
	// version -- os-arch
	goreleaserUrlFormat = goreleaserBaseRepo + "/releases/download/v%s/goreleaser_%s.%s"
)

type Goreleaser struct {
	version update.Version
	db      store.Storer
}

func (g *Goreleaser) IsService() bool {
	return false
}

func NewGoreleaser(db store.Storer) *Goreleaser {
	return &Goreleaser{
		db: db,
	}
}

func (g *Goreleaser) Version() update.Version {
	return g.version
}

func (g *Goreleaser) SetVersion(v update.Version) {
	g.version = v
}

func (g *Goreleaser) Name() string {
	return "goreleaser"
}

func (g *Goreleaser) Download() error {
	downloadDir := getDlPath(g.Name(), g.version.String())
	_ = os.MkdirAll(downloadDir, 0755)
	osname := fmt.Sprintf("%s_%s", osutil.GetOS(), osutil.GetAltArch())
	var ext string
	switch osutil.GetOS() {
	case "linux", "darwin":
		ext = "tar.gz"
	case "windows":
		ext = "zip"
	}
	fetchUrl := fmt.Sprintf(goreleaserUrlFormat, g.version, osname, ext)
	err := downloader.Download(fetchUrl, downloadDir)
	if err != nil {
		return err
	}
	return nil
}

func (g *Goreleaser) Install() error {
	var err error
	// install to global path
	// create bin directory under $PREFIX
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(g.Name(), g.version.String())
	// all files that are going to be installed
	executableName := g.Name()
	switch osutil.GetOS() {
	case "windows":
		executableName += ".exe"
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(binDir, executableName), 0755},
	}
	// copy file to the global bin directory
	ip, err := commonInstall(g, filesMap)
	if err != nil {
		return err
	}
	if err = g.db.New(ip); err != nil {
		return err
	}
	return nil
}

func (g *Goreleaser) Uninstall() error {
	var err error
	// install to global path
	// all files that are going to be installed
	dlPath := getDlPath(g.Name(), g.version.String())
	var pkg *store.InstalledPackage
	pkg, err = g.db.Get(g.Name())
	if err != nil {
		return err
	}
	var filesList []string
	for k := range pkg.FilesMap {
		filesList = append(filesList, k)
	}
	// uninstall listed files
	for _, file := range filesList {
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}
	// remove downloaded files
	return os.RemoveAll(dlPath)
}

func (g *Goreleaser) Update(version update.Version) error {
	g.version = version
	if err := g.Uninstall(); err != nil {
		return err
	}
	if err := g.Download(); err != nil {
		return err
	}
	return g.Install()
}

func (g *Goreleaser) Run(args ...string) error {
	pkg, err := g.db.Get(g.Name())
	if err != nil {
		return err
	}
	for k, _ := range pkg.FilesMap {
		if strings.Contains(k, g.Name()) {
			return osutil.Exec(k, args...)
		}
	}
	return nil
}

func (g *Goreleaser) Backup() error {
	// We don't need to implement this
	return nil
}

func (g *Goreleaser) RunStop() error {
	// We don't need to implement this
	return nil
}

func (g *Goreleaser) Dependencies() []dep.Component {
	return nil
}

func (g *Goreleaser) IsDev() bool {
	return true
}

func (g *Goreleaser) RepoUrl() update.RepositoryURL {
	return goreleaserBaseRepo
}
