package components

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
)

const (
	// version -- version -- os_alt_arch
	jsonnetUrl = "https://github.com/google/go-jsonnet"
)

type GoJsonnet struct {
	version update.Version
	db      store.Storer
}

func (g *GoJsonnet) IsService() bool {
	return false
}

func (g *GoJsonnet) Name() string {
	return "jsonnet"
}

func (g *GoJsonnet) Version() update.Version {
	return update.Version(g.version)
}

func (g *GoJsonnet) Download() error {
	targetDir := getDlPath(g.Name(), g.version.String())
	if osutil.DirExists(targetDir) {
		return downloader.GitCheckout("v"+g.version.String(), targetDir)
	}
	return downloader.GitClone(jsonnetUrl, targetDir, "v"+g.version.String())
}

func (g *GoJsonnet) Dependencies() []dep.Component {
	return nil
}

func (g *GoJsonnet) Install() error {
	var err error
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(g.Name(), g.version.String())
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
		filepath.Join(dlPath, g.Name()): {filepath.Join(binDir, g.Name()), 0755},
		//filepath.Join(dlPath, g.Name()+"fmt"): {filepath.Join(binDir, g.Name()+"fmt"), 0755},
	}

	ip, err := commonInstall(g, filesMap)
	if err != nil {
		return err
	}

	if err = g.db.New(ip); err != nil {
		return err
	}
	return nil
}

func (g *GoJsonnet) Uninstall() error {
	var err error
	err = os.RemoveAll(getDlPath(g.Name(), g.version.String()))
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

func (g *GoJsonnet) SetVersion(v update.Version) {
	g.version = v
}

func (g *GoJsonnet) Update(version update.Version) error {
	g.SetVersion(version)
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

func NewGoJsonnet(db store.Storer) *GoJsonnet {
	return &GoJsonnet{
		db: db,
	}
}

func (g *GoJsonnet) IsDev() bool {
	return true
}

func (g *GoJsonnet) RepoUrl() update.RepositoryURL {
	return jsonnetUrl
}
