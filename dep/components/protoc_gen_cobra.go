package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
)

const (
	// version -- version -- os_arch -- ext
	protoCobraRepo      = "https://github.com/amplify-edge/protoc-gen-cobra"
	protoCobraUrlFormat = protoCobraRepo + "/releases/download/v%s/protoc-gen-cobra-%s-%s.%s"
)

type ProtocGenCobra struct {
	version update.Version
	db      store.Storer
}

func (p *ProtocGenCobra) IsService() bool {
	return false
}

func NewProtocGenCobra(db store.Storer) *ProtocGenCobra {
	return &ProtocGenCobra{
		db: db,
	}
}

func (p *ProtocGenCobra) Name() string {
	return "protoc-gen-cobra"
}

func (p *ProtocGenCobra) Version() update.Version {
	if v := commonGetVersion(p, p.db); v != nil {
		return *v
	}
	return p.version
}

func (p *ProtocGenCobra) SetVersion(v update.Version) {
	p.version = v
}

func (p *ProtocGenCobra) Download() error {
	target := getDlPath(p.Name(), p.version.String())
	var ext string
	switch osutil.GetOS() {
	case "linux", "darwin":
		ext = "tar.gz"
	case "windows":
		ext = "zip"
	}
	osName := fmt.Sprintf("%s_%s", osutil.GetOS(), osutil.GetArch())
	fetchUrl := fmt.Sprintf(protoCobraUrlFormat, p.version, p.version, osName, ext)
	err := downloader.Download(fetchUrl, target)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenCobra) Install() error {
	var err error
	goBinDir := filepath.Join(osutil.GetGoPath(), "bin")
	_ = os.MkdirAll(goBinDir, 0755)

	executableName := p.Name()
	if osutil.GetOS() == "windows" {
		executableName += ".exe"
	}
	// all files that are going to be installed
	dlPath := getDlPath(p.Name(), p.version.String())
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(goBinDir, executableName), 0755},
	}

	ip, err := commonInstall(p, filesMap)
	if err != nil {
		return err
	}
	if err = p.db.New(ip); err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenCobra) Uninstall() error {
	var err error
	var pkg *store.InstalledPackage
	pkg, err = p.db.Get(p.Name())
	if err != nil {
		return err
	}
	var filesList []string
	for k := range pkg.FilesMap {
		filesList = append(filesList, k)
	}
	for _, f := range filesList {
		if err = os.RemoveAll(f); err != nil {
			return err
		}
	}
	if err = p.db.Delete(p.Name()); err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenCobra) Run(args ...string) error {
	return nil
}

func (p *ProtocGenCobra) Update(version update.Version) error {
	return commonUpdate(p, version)
}

func (p *ProtocGenCobra) RunStop() error {
	return nil
}

func (p *ProtocGenCobra) Backup() error {
	return nil
}

func (p *ProtocGenCobra) Dependencies() []dep.Component {
	return nil
}

func (p *ProtocGenCobra) IsDev() bool {
	return true
}

func (p *ProtocGenCobra) RepoUrl() update.RepositoryURL {
	return protoCobraRepo
}
