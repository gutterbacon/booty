package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"os"
	"path/filepath"
)

const (
	// version -- version -- os.arch
	protocGenGoUrlFormat = "https://github.com/protocolbuffers/protobuf-go/releases/download/v%s/protoc-gen-go.v%s.%s.tar.gz"
)

type ProtocGenGo struct {
	version string
	db      *store.DB
}

func (p *ProtocGenGo) Name() string {
	return "protoc-gen-go"
}

func (p *ProtocGenGo) Version() string {
	return p.version
}

func (p *ProtocGenGo) Download() error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	osName := fmt.Sprintf("%s.%s", osutil.GetOS(), osutil.GetArch())
	fetchUrl := fmt.Sprintf(protocGenGoUrlFormat, p.version, p.version, osName)
	target := getDlPath(p.Name(), p.version)

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
	dlPath := getDlPath(p.Name(), p.version)
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(goBinDir, executableName), 0755},
	}

	ip := store.InstalledPackage{
		Name:     p.Name(),
		Version:  p.version,
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
	if err = p.db.New(&ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
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
	return nil
}

func (p *ProtocGenGo) Run(args ...string) error {
	// Should not run as standalone, rather to compliment protoc
	return nil
}

func (p *ProtocGenGo) Update(version string) error {
	p.version = version
	if err := p.Uninstall(); err != nil {
		return err
	}
	if err := p.Download(); err != nil {
		return err
	}
	return p.Install()
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

func NewProtocGenGo(db *store.DB, version string) *ProtocGenGo {
	return &ProtocGenGo{
		version: version,
		db:      db,
	}
}