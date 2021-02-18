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
	"go.amplifyedge.org/booty-v2/internal/fileutil"

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
	ip := store.InstalledPackage{
		Name:     g.Name(),
		Version:  g.version.String(),
		FilesMap: map[string]int{},
	}
	// copy file to the global bin directory
	for k, v := range filesMap {
		if err = fileutil.Copy(k, v[0].(string)); err != nil {
			return err
		}
		installedName := v[0].(string)
		installedMode := v[1].(int)
		if err = os.Chmod(installedName, os.FileMode(installedMode)); err != nil {
			return err
		}
		ip.FilesMap[installedName] = installedMode
	}
	if err = g.db.New(&ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
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
