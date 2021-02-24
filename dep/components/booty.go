package components

import (
	"fmt"
	ks "github.com/kardianos/service"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/service"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
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
	executablePath := filepath.Join(osutil.GetGoPath(), "bin", ex)
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
		bootyRepoUrl+"/releases/download/v%s/booty-%s-%s_%s.%s",
		ver, ver, osutil.GetOS(), osutil.GetArch(), ext,
	)
	return downloader.Download(fetchUrl, getDlPath(b.Name(), ver))
}

func (b *Booty) Dependencies() []dep.Component {
	return nil
}

func (b *Booty) Install() error {
	goBinDir := filepath.Join(osutil.GetGoPath(), "bin")
	dlPath := getDlPath(b.Name(), b.version.String())
	executableName := b.Name()
	if osutil.GetOS() == "windows" {
		executableName += ".exe"
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(goBinDir, executableName)},
	}
	ip, err := commonInstall(b, filesMap)
	if err != nil {
		return err
	}
	_ = b.svc.Install()
	return b.db.New(ip)
}

func (b *Booty) Uninstall() error {
	// can't uninstall booty
	return nil
}

func (b *Booty) Run(args ...string) error {
	return b.svc.Start()
}

func (b *Booty) Update(version update.Version) error {
	return commonUpdate(b, b.version)
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
