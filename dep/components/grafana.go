package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/pkg/fileutil"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"

	"go.amplifyedge.org/booty-v2/pkg/downloader"
)

const (
	// version -- os-arch
	fetchUrlFormat = "https://dl.grafana.com/oss/release/grafana-%s.%s.%s"
)

// Grafana implements Component interface
type Grafana struct {
	version string
	dlPath  string
	db      *store.DB
}

func NewGrafana(db *store.DB, version string) *Grafana {
	return &Grafana{version, "", db}
}

// Gets grafana's version
func (g *Grafana) Version() string {
	return g.version
}

func (g *Grafana) Name() string {
	return "grafana"
}

func (g *Grafana) Download(targetDir string) error {
	osname := fmt.Sprintf("%s-%s", strings.ToLower(osutil.GetOS()), osutil.GetArch())
	var fetchUrl string
	switch osutil.GetOS() {
	case "linux", "darwin":
		fetchUrl = fmt.Sprintf(fetchUrlFormat, g.version, osname, "tar.gz")
	case "windows":
		fetchUrl = fmt.Sprintf(fetchUrlFormat, g.version, osname, "zip")
	}
	err := downloader.Download(fetchUrl, targetDir)
	if err != nil {
		return err
	}
	g.dlPath = filepath.Join(targetDir, g.Name()+"-"+g.version)
	return nil
}

func (g *Grafana) Install() error {
	var err error
	// install to path
	binDir := osutil.GetBinDir()
	etcDir := osutil.GetEtcDir()

	serverExecutable := g.Name() + "-server"
	clientExecutable := g.Name() + "-cli"

	switch osutil.GetOS() {
	case "windows":
		serverExecutable += ".exe"
		clientExecutable += ".exe"
	}

	// all files that are going to be installed
	filesMap := map[string][]interface{}{
		filepath.Join(g.dlPath, "bin", serverExecutable): {filepath.Join(binDir, serverExecutable), 0755},
		filepath.Join(g.dlPath, "bin", clientExecutable): {filepath.Join(binDir, clientExecutable), 0755},
		filepath.Join(g.dlPath, "conf", "defaults.ini"):  {filepath.Join(etcDir, "grafana.ini"), 0644},
		filepath.Join(g.dlPath, "conf", "sample.ini"):    {filepath.Join(etcDir, "grafana.sample.ini"), 0644},
	}

	ip := store.InstalledPackage{
		Name:     g.Name(),
		Version:  g.version,
		FilesMap: map[string]int{},
	}

	// copy file to the bin directory
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
	return os.RemoveAll(g.dlPath)
}

func (g *Grafana) Uninstall() error {
	var err error
	// install to global path
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
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}
	// remove downloaded files
	return os.RemoveAll(g.dlPath)
}

func (g *Grafana) Update(version string) error {
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

func (g *Grafana) Run(args ...string) error {
	return nil
}

func (g *Grafana) Backup() error {
	return nil
}

func (g *Grafana) RunStop() error {
	return nil
}
