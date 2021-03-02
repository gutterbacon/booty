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
	protocGenGoRepo = "https://github.com/protocolbuffers/protobuf-go"
	// version -- version -- os.arch
	protocGenGoUrlFormat = protocGenGoRepo + "/releases/download/v%s/protoc-gen-go.v%s.%s.tar.gz"
)

type ProtocGenGo struct {
	version update.Version
	db      store.Storer
}

func (p *ProtocGenGo) IsService() bool {
	return false
}

func (p *ProtocGenGo) Name() string {
	return "protoc-gen-go"
}

func (p *ProtocGenGo) Version() update.Version {
	if v := commonGetVersion(p, p.db); v != nil {
		return *v
	}
	return p.version
}

func (p *ProtocGenGo) SetVersion(v update.Version) {
	p.version = v
}

func (p *ProtocGenGo) Download() error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	osName := fmt.Sprintf("%s.%s", osutil.GetOS(), osutil.GetArch())
	fetchUrl := fmt.Sprintf(protocGenGoUrlFormat, p.version, p.version, osName)
	target := getDlPath(p.Name(), p.version.String())

	err := downloader.Download(fetchUrl, target)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenGo) Install() error {
	var err error
	// install to path
	goBinDir := filepath.Join(osutil.GetGoPath(), "bin")
	_ = os.MkdirAll(goBinDir, 0755)

	executableName := p.Name()

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

func (p *ProtocGenGo) Uninstall() error {
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
	for _, file := range filesList {
		if err = os.RemoveAll(file); err != nil {
			return err
		}
	}
	if err = p.db.Delete(p.Name()); err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenGo) Run(args ...string) error {
	// Should not run as standalone, rather to compliment protoc
	return nil
}

func (p *ProtocGenGo) Update(version update.Version) error {
	return commonUpdate(p, version)
}

func (p *ProtocGenGo) RunStop() error {
	return nil
}

func (p *ProtocGenGo) Backup() error {
	return nil
}

func (p *ProtocGenGo) Dependencies() []dep.Component {
	return nil
}

func NewProtocGenGo(db store.Storer) *ProtocGenGo {
	return &ProtocGenGo{
		db: db,
	}
}

func (p *ProtocGenGo) IsDev() bool {
	return true
}

func (p *ProtocGenGo) RepoUrl() update.RepositoryURL {
	return protocGenGoRepo
}
