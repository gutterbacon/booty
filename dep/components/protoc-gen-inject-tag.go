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
	injectRepo   = "github.com/favadi/protoc-go-inject-tag"
	injectUrlFmt = "https://" + injectRepo
)

type ProtocGenInjectTag struct {
	version update.Version
	db      store.Storer
}

func NewProtocGenInjectTag(db store.Storer) *ProtocGenInjectTag {
	return &ProtocGenInjectTag{db: db}
}

func (p *ProtocGenInjectTag) Name() string {
	return "protoc-gen-inject-tag"
}

func (p *ProtocGenInjectTag) Version() update.Version {
	return p.version
}

func (p *ProtocGenInjectTag) SetVersion(version update.Version) {
	p.version = version
}

func (p *ProtocGenInjectTag) Download() error {
	targetDir := getDlPath(p.Name(), p.version.String())
	return downloader.Download(injectRepo+"?ref=v"+p.version.String(), targetDir)
}

func (p *ProtocGenInjectTag) Dependencies() []dep.Component {
	return nil
}

func (p *ProtocGenInjectTag) Install() error {
	dlDir := getDlPath(p.Name(), p.version.String())
	var err error
	goBinDir := filepath.Join(osutil.GetGoPath(), "bin")
	_ = os.MkdirAll(goBinDir, 0755)
	err = os.Chdir(dlDir)
	if err != nil {
		return err
	}
	if err = osutil.Exec("go", "build", "-ldflags", `-s -w`, "-o", p.Name(), "."); err != nil {
		return err
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlDir, p.Name()): {filepath.Join(goBinDir, p.Name()), 0755},
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

func (p *ProtocGenInjectTag) Uninstall() error {
	var err error
	err = os.RemoveAll(getDlPath(p.Name(), p.version.String()))
	if err != nil {
		return err
	}
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

func (p *ProtocGenInjectTag) Run(args ...string) error {
	return osutil.Exec(p.Name(), args...)
}

func (p *ProtocGenInjectTag) Update(version update.Version) error {
	return commonUpdate(p, version)
}

func (p *ProtocGenInjectTag) RunStop() error {
	return nil
}

func (p *ProtocGenInjectTag) Backup() error {
	return nil
}

func (p *ProtocGenInjectTag) IsDev() bool {
	return true
}

func (p *ProtocGenInjectTag) IsService() bool {
	return false
}

func (p *ProtocGenInjectTag) RepoUrl() update.RepositoryURL {
	return injectUrlFmt
}
