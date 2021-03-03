package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
	"strings"
)

const (
	hugoBase = "github.com/gohugoio/hugo"
	// -- version -- version -- UpperCase OS - Alt Arch -- Extension
	hugoUrlFmt = "https://" + hugoBase + "/releases/download/v%s/hugo_%s_%s-%s.%s"
)

type Hugo struct {
	version update.Version
	db      store.Storer
}

func NewHugo(db store.Storer) *Hugo {
	return &Hugo{
		db: db,
	}
}

func (h *Hugo) Name() string {
	return "hugo"
}

func (h *Hugo) Version() update.Version {
	if v := commonGetVersion(h, h.db); v != nil {
		return *v
	}
	return h.version
}

func (h *Hugo) SetVersion(version update.Version) {
	h.version = version
}

func (h *Hugo) Download() error {
	dlDir := getDlPath(h.Name(), h.Version().String())
	_ = os.MkdirAll(dlDir, 0755)
	var osname string
	var ext string
	var arch string
	switch osutil.GetOS() {
	case "linux":
		osname = strings.ToUpper(osutil.GetOS())
		ext = "tar.gz"
	case "darwin":
		osname = "macOS"
		ext = "tar.gz"
	case "windows":
		osname = strings.ToUpper(osutil.GetOS())
		ext = "zip"
	}
	switch osutil.GetArch() {
	case "amd64":
		arch = "64bit"
	case "aarch64":
		arch = "ARM64"
	}
	fetchUrl := fmt.Sprintf(hugoUrlFmt, h.version, h.version, osname, arch, ext)
	return downloader.Download(fetchUrl, dlDir)
}

func (h *Hugo) Dependencies() []dep.Component {
	return nil
}

func (h *Hugo) Install() error {
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(h.Name(), h.version.String())
	executableName := h.Name()
	switch osutil.GetOS() {
	case "windows":
		executableName += ".exe"
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(binDir, executableName), 0755},
	}
	ip, err := commonInstall(h, filesMap)
	if err != nil {
		return err
	}
	return h.db.New(ip)
}

func (h *Hugo) Uninstall() error {
	var err error
	dlPath := getDlPath(h.Name(), h.version.String())
	var pkg *store.InstalledPackage
	pkg, err = h.db.Get(h.Name())
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
	if err = h.db.Delete(h.Name()); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
}

func (h *Hugo) Run(args ...string) error {
	return osutil.Exec(filepath.Join(osutil.GetBinDir(), h.Name()), args...)
}

func (h *Hugo) Update(version update.Version) error {
	return commonUpdate(h, version)
}

func (h *Hugo) RunStop() error {
	return nil
}

func (h *Hugo) Backup() error {
	return nil
}

func (h *Hugo) IsDev() bool {
	return false
}

func (h *Hugo) IsService() bool {
	return false
}

func (h *Hugo) RepoUrl() update.RepositoryURL {
	return "https://" + hugoBase
}
