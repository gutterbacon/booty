package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/pkg/downloader"
	"go.amplifyedge.org/booty-v2/pkg/fileutil"
	"os"
	"path/filepath"
	"strings"

	//"path/filepath"

	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"
)

const (
	// version -- os-arch
	goreleaserUrlFormat = "https://github.com/goreleaser/goreleaser/releases/download/v%s/goreleaser_%s.tar.gz"
)

type Goreleaser struct {
	version string
	dlPath  string
	db      *store.DB
}

func NewGoreleaser(db *store.DB, version string) *Goreleaser {
	return &Goreleaser{
		version: version,
		dlPath:  "",
		db:      db,
	}
}

func (g *Goreleaser) Version() string {
	return g.version
}

func (g *Goreleaser) Name() string {
	return "goreleaser"
}

func (g *Goreleaser) Download(targetDir string) error {
	downloadDir := filepath.Join(targetDir, g.Name())
	_ = os.MkdirAll(downloadDir, 0755)
	osname := fmt.Sprintf("%s_%s", osutil.GetOS(), osutil.GetAltArch())
	fetchUrl := fmt.Sprintf(goreleaserUrlFormat, g.version, osname)
	err := downloader.Download(fetchUrl, downloadDir)
	if err != nil {
		return err
	}
	g.dlPath = downloadDir
	return nil
}

func (g *Goreleaser) Install() error {
	var err error
	// install to global path
	switch strings.ToLower(osutil.GetOS()) {
	case "linux", "darwin":
		// create bin directory under $PREFIX
		binDir := osutil.GetBinDir()
		// all files that are going to be installed
		filesMap := map[string][]interface{}{
			filepath.Join(g.dlPath, "goreleaser"): {filepath.Join(binDir, "goreleaser"), 0755},
		}
		ip := store.InstalledPackage{
			Name:     g.Name(),
			Version:  g.version,
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
	}
	return os.RemoveAll(g.dlPath)
}

func (g *Goreleaser) Uninstall() error {
	var err error
	// install to global path
	switch strings.ToLower(osutil.GetOS()) {
	case "linux", "darwin":
		// all files that are going to be installed
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
			if err = osutil.ExecSudo("rm", "-rf", file); err != nil {
				return err
			}
		}
	}
	// remove downloaded files
	return os.RemoveAll(g.dlPath)
}

func (g *Goreleaser) Update(version string) error {
	g.version = version
	targetDir := filepath.Dir(g.dlPath)
	if err := g.Uninstall(); err != nil {
		return err
	}
	if err := g.Download(targetDir); err != nil {
		return err
	}
	return g.Install()
}

func (g *Goreleaser) Run() error {
	pkg, err := g.db.Get(g.Name())
	if err != nil {
		return err
	}
	for k := range pkg.FilesMap {
		if k == g.Name() {
			return osutil.Exec(k)
		}
	}
	return nil
}

func (g *Goreleaser) Backup() error {
	return nil
}

func (g *Goreleaser) Stop() error {
	return nil
}