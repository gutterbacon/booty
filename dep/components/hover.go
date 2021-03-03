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
	// https://github.com/go-flutter-desktop/hover/archive/v0.46.2.zip
	hoverBase = "github.com/go-flutter-desktop/hover"
)

type Hover struct {
	db      store.Storer
	version update.Version
}

func NewHover(db store.Storer) *Hover {
	return &Hover{
		db: db,
	}
}

func (h *Hover) Name() string {
	return "hover"
}

func (h *Hover) Version() update.Version {
	if v := commonGetVersion(h, h.db); v != nil {
		return *v
	}
	return h.version
}

func (h *Hover) SetVersion(version update.Version) {
	h.version = version
}

func (h *Hover) Download() error {
	dlDir := getDlPath(h.Name(), h.version.String())
	return downloader.Download(hoverBase+"?ref=v"+h.version.String(), dlDir)
}

func (h *Hover) Dependencies() []dep.Component {
	return nil
}

func (h *Hover) Install() error {
	var err error
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(h.Name(), h.version.String())
	// change dir to download path
	if err = os.Chdir(dlPath); err != nil {
		return err
	}
	// build binaries
	if err = osutil.Exec("go", "mod", "tidy"); err != nil {
		return err
	}
	if err = osutil.Exec("go", "build", "-ldflags", `-s -w`, "-o", h.Name(), "."); err != nil {
		return nil
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, h.Name()): {filepath.Join(binDir, h.Name()), 0755},
	}
	ip, err := commonInstall(h, filesMap)
	if err != nil {
		return err
	}
	return h.db.New(ip)
}

func (h *Hover) Uninstall() error {
	var err error
	var pkg *store.InstalledPackage
	pkg, err = h.db.Get(h.Name())
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
	if err = h.db.Delete(h.Name()); err != nil {
		return err
	}
	return os.RemoveAll(getDlPath(h.Name(), h.version.String()))
}

func (h *Hover) Run(args ...string) error {
	return osutil.Exec(filepath.Join(osutil.GetBinDir(), h.Name()), args...)
}

func (h *Hover) Update(version update.Version) error {
	return commonUpdate(h, version)
}

func (h *Hover) RunStop() error {
	return nil
}

func (h *Hover) Backup() error {
	return nil
}

func (h *Hover) IsDev() bool {
	return true
}

func (h *Hover) IsService() bool {
	return false
}

func (h *Hover) RepoUrl() update.RepositoryURL {
	return "https://"+hoverBase
}
