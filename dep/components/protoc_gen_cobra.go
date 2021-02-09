package components

import (
	"go.amplifyedge.org/booty-v2/pkg/downloader"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"go.amplifyedge.org/booty-v2/pkg/store"
	"os"
	"path/filepath"
)

const (
	protoCobraUrl = "https://github.com/amplify-edge/protoc-gen-cobra/archive/v0.4.0.tar.gz"
)

type ProtocGenCobra struct {
	version string
	dlPath  string
	db      *store.DB
}

func NewProtocGenCobra(db *store.DB, version string) *ProtocGenCobra {
	return &ProtocGenCobra{
		version: version,
		dlPath:  "",
		db:      db,
	}
}

func (p *ProtocGenCobra) Name() string {
	return "protoc-gen-cobra"
}

func (p *ProtocGenCobra) Version() string {
	return p.version
}

func (p *ProtocGenCobra) Download(targetDir string) error {
	target := filepath.Join(targetDir, p.Name()+"-"+p.version)
	err := downloader.Download(protoCobraUrl, targetDir)
	if err != nil {
		return err
	}
	p.dlPath = target
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
	_ = os.Chdir(p.dlPath)
	if err = osutil.Exec("go", "install"); err != nil {
		return err
	}
	ip := store.InstalledPackage{
		Name:    p.Name(),
		Version: p.version,
		FilesMap: map[string]int{
			filepath.Join(goBinDir, executableName): 0755,
		},
	}
	if err = p.db.New(&ip); err != nil {
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
	return nil
}

func (p *ProtocGenCobra) Run(args ...string) error {
	return nil
}

func (p *ProtocGenCobra) Update(version string) error {
	p.version = version
	targetDir := filepath.Dir(p.dlPath)
	if err := p.Uninstall(); err != nil {
		return err
	}
	if err := p.Download(targetDir); err != nil {
		return err
	}
	return p.Install()
}

func (p *ProtocGenCobra) RunStop() error {
	return nil
}

func (p *ProtocGenCobra) Backup() error {
	return nil
}
