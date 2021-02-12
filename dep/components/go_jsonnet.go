package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
)

const (
	// version -- version -- os_alt_arch
	jsonnetUrlFmt = "https://github.com/google/go-jsonnet/releases/download/v%s/go-jsonnet_%s_%s.tar.gz"
)

type GoJsonnet struct {
	version string
	db *store.DB
}

func (g *GoJsonnet) Name() string {
	return "jsonnet"
}

func (g *GoJsonnet) Version() string {
	return g.version
}

func (g *GoJsonnet) Download() error {
	targetDir := getDlPath(g.Name(), g.version)
	if osutil.GetOS() == "darwin" {
		_ = os.Setenv("GOPATH", targetDir)
		_ = os.Setenv("GO111MODULES", "off")
		return osutil.Exec("go", "get", "-u", "-v", "github.com/google/go-jsonnet/cmd/jsonnet@v" + g.version)
	}
	var osVer string
	if osutil.GetOS() == "linux" && osutil.GetArch() == "arm64" {
		osVer = "Linux_arm64"
	} else {
		osVer = fmt.Sprintf("%s_%s", strings.ToUpper(osutil.GetOS()), osutil.GetAltArch())
	}
	fetchUrl := fmt.Sprintf(jsonnetUrlFmt, g.version, g.version, osVer)
	return downloader.Download(fetchUrl, targetDir)
}

func (g *GoJsonnet) Dependencies() []dep.Component {
	return nil
}

func (g *GoJsonnet) Install() error {
	var err error
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(g.Name(), g.version)
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, "bin", g.Name()): {filepath.Join(binDir, g.Name()), 0755},
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
	return os.RemoveAll(dlPath)
}

func (g *GoJsonnet) Uninstall() error {
	var err error
	var pkg *store.InstalledPackage
	pkg, err = g.db.Get(g.Name())
	if err != nil {
		return err
	}
	var filesList []string
	for k := range pkg.FilesMap {
		filesList = append(filesList, k)
	}
	for _, file := range filesList {
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}
	return nil
}

func (g *GoJsonnet) Run(args ...string) error {
	panic("implement me")
}

func (g *GoJsonnet) Update(version string) error {
	g.version = version
	if err := g.Uninstall(); err != nil {
		return err
	}
	if err := g.Download(); err != nil {
		return err
	}
	return g.Install()
}

func (g *GoJsonnet) RunStop() error {
	return nil
}

func (g *GoJsonnet) Backup() error {
	return nil
}

func NewGoJsonnet(db *store.DB, version string) *GoJsonnet {
	return &GoJsonnet{
		version: version,
		db:      db,
	}
}

