package components

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"os"
	"path/filepath"
)

const (
	// version -- version -- os_alt_arch
	jsonnetUrl = "https://github.com/google/go-jsonnet.git"
)

type GoJsonnet struct {
	version string
	db      *store.DB
}

func (g *GoJsonnet) Name() string {
	return "jsonnet"
}

func (g *GoJsonnet) Version() string {
	return g.version
}

func (g *GoJsonnet) Download() error {
	targetDir := getDlPath(g.Name(), g.version)
	if osutil.DirExists(targetDir) {
		return nil
	}
	return downloader.GitClone(jsonnetUrl, targetDir, "v"+g.version)
}

func (g *GoJsonnet) Dependencies() []dep.Component {
	return nil
}

func (g *GoJsonnet) Install() error {
	var err error
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(g.Name(), g.version)
	// change dir to download path
	if err = os.Chdir(dlPath); err != nil {
		return err
	}
	// build binaries
	recipes := []string{g.Name(), g.Name() + "fmt"}
	for _, r := range recipes {
		if err = osutil.Exec("go", "build", "-o", r, filepath.Join(dlPath, "cmd", r)); err != nil {
			return err
		}
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, g.Name()):       {filepath.Join(binDir, g.Name()), 0755},
		filepath.Join(dlPath, g.Name()+"fmt"): {filepath.Join(binDir, g.Name()+"fmt"), 0755},
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
	return nil
}

func (g *GoJsonnet) Uninstall() error {
	var err error
	err = os.RemoveAll(getDlPath(g.Name(), g.version))
	if err != nil {
		return err
	}
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
	return osutil.Exec(filepath.Join(osutil.GetBinDir(), g.Name()), args...)
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
