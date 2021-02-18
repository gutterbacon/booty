package components

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"

	"fmt"
	"os"
	"path/filepath"
)

const (
	// version -- version -- os.arch
	genRpcRepo       = "https://github.com/grpc/grpc-go"
	genGrpcUrlFormat = genRpcRepo + "/releases/download/cmd/protoc-gen-go-grpc/v%s/protoc-gen-go-grpc.v%s.%s.tar.gz"
)

type ProtocGenGoGrpc struct {
	version update.Version
	db      store.Storer
}

func (p *ProtocGenGoGrpc) IsService() bool {
	return false
}

func (p *ProtocGenGoGrpc) Name() string {
	return "protoc-gen-go-grpc"
}

func (p *ProtocGenGoGrpc) Version() update.Version {
	return p.version
}

func (p *ProtocGenGoGrpc) SetVersion(v update.Version) {
	p.version = v
}

func (p *ProtocGenGoGrpc) Download() error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	osName := fmt.Sprintf("%s.%s", osutil.GetOS(), osutil.GetArch())
	fetchUrl := fmt.Sprintf(genGrpcUrlFormat, p.version, p.version, osName)
	target := getDlPath(p.Name(), p.version.String())
	err := downloader.Download(fetchUrl, target)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProtocGenGoGrpc) Install() error {
	var err error
	// install to path
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
	return os.RemoveAll(dlPath)
}

func (p *ProtocGenGoGrpc) Uninstall() error {
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

func (p *ProtocGenGoGrpc) Run(args ...string) error {
	// Should not run as standalone, rather to compliment protoc
	return nil
}

func (p *ProtocGenGoGrpc) Update(version update.Version) error {
	p.version = version
	if err := p.Uninstall(); err != nil {
		return err
	}
	if err := p.Download(); err != nil {
		return err
	}
	return p.Install()
}

func (p *ProtocGenGoGrpc) RunStop() error {
	return nil
}

func (p *ProtocGenGoGrpc) Backup() error {
	return nil
}

func (p *ProtocGenGoGrpc) Dependencies() []dep.Component {
	return nil
}

func (p *ProtocGenGoGrpc) IsDev() bool {
	return true
}

func (p *ProtocGenGoGrpc) RepoUrl() update.RepositoryURL {
	return genRpcRepo
}

func NewProtocGenGoGrpc(db store.Storer) *ProtocGenGoGrpc {
	return &ProtocGenGoGrpc{
		db: db,
	}
}
