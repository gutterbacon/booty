package components

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	gupd "github.com/inconshreveable/go-update"
	"github.com/kardianos/osext"
	ks "github.com/kardianos/service"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/service"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Booty struct {
	version update.Version
	db      store.Storer
	svc     *service.Svc
}

const (
	bootyRepoBase = "github.com/amplify-edge/booty"
	bootyRepoUrl  = "https://" + bootyRepoBase
)

func NewBooty(db store.Storer) *Booty {
	booty := &Booty{
		db: db,
	}
	svc, err := booty.service()
	if err != nil {
		fmt.Println("failed to bootstrap booty")
		os.Exit(1)
	}
	booty.svc = svc
	return booty
}

func (b *Booty) service() (*service.Svc, error) {
	ex := b.Name()
	if osutil.GetOS() == "windows" {
		ex += ".exe"
	}
	// lookup for absolute booty path
	executablePath, err := exec.LookPath(ex)
	if err != nil {
		// if it's not found get this program directory
		executablePath, err = osext.ExecutableFolder()
		if err != nil {
			return nil, err
		}
	}
	config := &ks.Config{
		Name:        b.Name(),
		DisplayName: b.Name(),
		Description: "bootstraps amplify-edge projects along with its dependencies",
		Arguments:   []string{"agent"},
		Executable:  executablePath,
		Option:      map[string]interface{}{},
	}
	return service.NewService(config)
}

func (b *Booty) Name() string {
	return "booty"
}

func (b *Booty) Version() update.Version {
	if v := commonGetVersion(b, b.db); v != nil {
		return *v
	}
	return b.version
}

func (b *Booty) SetVersion(version update.Version) {
	b.version = version
}

func (b *Booty) Download() error {
	var ext string
	switch osutil.GetOS() {
	case "linux", "darwin":
		ext = "tar.gz"
	case "windows":
		ext = "zip"
	}
	ver := b.version.String()
	fetchUrl := fmt.Sprintf(
		bootyRepoUrl+"/releases/download/v%s/booty-v%s-%s_%s.%s",
		ver, ver, osutil.GetOS(), osutil.GetArch(), ext,
	)
	return downloader.Download(fetchUrl, getDlPath(b.Name(), ver))
}

func (b *Booty) Dependencies() []dep.Component {
	return nil
}

func (b *Booty) Install() error {
	bootyExePath, err := osext.Executable()
	if err != nil {
		return err
	}
	dlPath := getDlPath(b.Name(), b.version.String())
	executableName := b.Name()
	if osutil.GetOS() == "windows" {
		executableName += ".exe"
	}

	content, err := ioutil.ReadFile(filepath.Join(dlPath, executableName))
	if err != nil {
		return err
	}

	ip := &store.InstalledPackage{
		Name:    b.Name(),
		Version: b.Version().String(),
		FilesMap: map[string]string{
			bootyExePath: fmt.Sprintf("%x", sha256.Sum256(content)),
		},
	}

	if err = gupd.Apply(bytes.NewBuffer(content), gupd.Options{TargetPath: bootyExePath}); err != nil {
		return err
	}
	if err = b.db.New(ip); err != nil {
		return err
	}
	_ = b.svc.Install()
	return nil
}

func (b *Booty) Uninstall() error {
	// can't uninstall booty
	return nil
}

func (b *Booty) Run(args ...string) error {
	return b.svc.Start()
}

func (b *Booty) Update(version update.Version) error {
	b.version = version
	return b.Install()
}

func (b *Booty) RunStop() error {
	return b.svc.Stop()
}

func (b *Booty) Backup() error {
	return nil
}

func (b *Booty) IsDev() bool {
	return false
}

func (b *Booty) IsService() bool {
	return true
}

func (b *Booty) RepoUrl() update.RepositoryURL {
	return bootyRepoUrl
}
