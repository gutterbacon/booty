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
)

const (
	mkcertRepo = "https://github.com/FiloSottile/mkcert"
	// version -- version -- os-arch (windows add exe)
	mkcertUrlFmt = mkcertRepo + "/releases/download/v%s/mkcert-v%s-%s"
)

type Mkcert struct {
	version update.Version
	db      store.Storer
	osname  string
}

func (m *Mkcert) Name() string {
	return "mkcert"
}

func (m *Mkcert) Version() update.Version {
	if v := commonGetVersion(m, m.db); v != nil {
		return *v
	}
	return m.version
}

func (m *Mkcert) SetVersion(version update.Version) {
	m.version = version
}

func (m *Mkcert) Download() error {
	dlDir := getDlPath(m.Name(), m.version.String())
	_ = os.MkdirAll(dlDir, 0755)

	fetchUrl := fmt.Sprintf(mkcertUrlFmt, m.version.String(), m.version.String(), m.osname)
	return downloader.Download(fetchUrl, dlDir)
}

func (m *Mkcert) Dependencies() []dep.Component {
	return nil
}

func (m *Mkcert) Install() error {
	var err error
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(m.Name(), m.version.String())
	executable := fmt.Sprintf("%s-v%s-%s", m.Name(), m.version.String(), m.osname)
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executable): {filepath.Join(binDir, m.Name()), 0755},
	}
	ip, err := commonInstall(m, filesMap)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(osutil.GetDataDir(), m.Name()), 0700)
	if err = m.db.New(ip); err != nil {
		return err
	}
	return nil
}

func (m *Mkcert) Uninstall() error {
	var err error
	// install to global path
	// all files that are going to be installed
	dlPath := getDlPath(m.Name(), m.version.String())
	var pkg *store.InstalledPackage
	pkg, err = m.db.Get(m.Name())
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
	if err = m.db.Delete(m.Name()); err != nil {
		return err
	}
	// remove downloaded files
	return os.RemoveAll(dlPath)
}

func (m *Mkcert) Run(args ...string) error {
	_ = os.Setenv("CAROOT", filepath.Join(osutil.GetDataDir(), m.Name()))
	return osutil.Exec(filepath.Join(osutil.GetBinDir(), m.Name()), args...)
}

func (m *Mkcert) Update(version update.Version) error {
	return commonUpdate(m, version)
}

func (m *Mkcert) RunStop() error {
	return nil
}

func (m *Mkcert) Backup() error {
	return nil
}

func (m *Mkcert) IsDev() bool {
	return true
}

func (m *Mkcert) IsService() bool {
	return false
}

func (m *Mkcert) RepoUrl() update.RepositoryURL {
	return mkcertRepo
}

func NewMkcert(db store.Storer) *Mkcert {
	osname := fmt.Sprintf("%s-%s", osutil.GetOS(), osutil.GetArch())
	if osname == "windows-amd64" {
		osname += ".exe"
	}
	return &Mkcert{db: db, osname: osname}
}
