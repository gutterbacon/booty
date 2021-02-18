package components

import (
	"fmt"
	ks "github.com/kardianos/service"
	"go.amplifyedge.org/booty-v2/internal/service"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
)

const (
	// version -- os-arch
	grafanaBaseRepo = "https://github.com/grafana/grafana"
	fetchUrlFormat  = "https://dl.grafana.com/oss/release/grafana-%s.%s.%s"
)

// Grafana implements Component interface
type Grafana struct {
	version update.Version
	db      store.Storer
	svc     *service.Svc
}

func (g *Grafana) IsService() bool {
	return true
}

func NewGrafana(db store.Storer) *Grafana {
	return &Grafana{db: db}
}

// Gets grafana's version
func (g *Grafana) Version() update.Version {
	return update.Version(g.version)
}

func (g *Grafana) SetVersion(v update.Version) {
	g.version = v
}

func (g *Grafana) Name() string {
	return "grafana"
}

func (g *Grafana) Download() error {
	osname := fmt.Sprintf("%s-%s", strings.ToLower(osutil.GetOS()), osutil.GetArch())
	var fetchUrl string
	switch osutil.GetOS() {
	case "linux", "darwin":
		fetchUrl = fmt.Sprintf(fetchUrlFormat, g.version, osname, "tar.gz")
	case "windows":
		fetchUrl = fmt.Sprintf(fetchUrlFormat, g.version, osname, "zip")
	}
	targetDir := osutil.GetDownloadDir()
	err := downloader.Download(fetchUrl, targetDir)
	if err != nil {
		return err
	}
	return nil
}

func (g *Grafana) service() (*service.Svc, error) {
	config := &ks.Config{
		Name:             g.Name(),
		DisplayName:      g.Name(),
		Description:      "Extensible platform that uses TLS by default",
		WorkingDirectory: filepath.Join(osutil.GetEtcDir(), "grafana"),
		Arguments: []string{
			"--config",
			filepath.Join(osutil.GetEtcDir(), "grafana", "grafana.ini"),
		},
		Executable: filepath.Join(osutil.GetBinDir(), g.Name()+"-server"),
		Option:     map[string]interface{}{},
	}
	return service.NewService(config)
}

func (g *Grafana) Install() error {
	var err error
	// install to path
	binDir := osutil.GetBinDir()
	grafanaEtcDir := filepath.Join(osutil.GetEtcDir(), "grafana")

	serverExecutable := g.Name() + "-server"
	clientExecutable := g.Name() + "-cli"

	switch osutil.GetOS() {
	case "windows":
		serverExecutable += ".exe"
		clientExecutable += ".exe"
	}
	dlPath := getDlPath(g.Name(), g.version.String())

	err = os.MkdirAll(grafanaEtcDir, 0755)

	// all files that are going to be installed
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, "bin", serverExecutable): {filepath.Join(binDir, serverExecutable), 0755},
		filepath.Join(dlPath, "bin", clientExecutable): {filepath.Join(binDir, clientExecutable), 0755},
		filepath.Join(dlPath, "conf", "defaults.ini"):  {filepath.Join(grafanaEtcDir, "grafana.ini"), 0644},
		filepath.Join(dlPath, "conf", "sample.ini"):    {filepath.Join(grafanaEtcDir, "grafana.sample.ini"), 0644},
		filepath.Join(dlPath, "conf"):                  {filepath.Join(grafanaEtcDir, "conf"), 0755},
		filepath.Join(dlPath, "plugins-bundled"):       {filepath.Join(grafanaEtcDir, "plugins-bundled"), 0755},
		filepath.Join(dlPath, "public"):                {filepath.Join(grafanaEtcDir, "public"), 0755},
		filepath.Join(dlPath, "scripts"):               {filepath.Join(grafanaEtcDir, "scripts"), 0755},
	}

	ip, err := commonInstall(g, filesMap)
	if err != nil {
		return err
	}

	// install service
	s, err := g.service()
	if err != nil {
		return err
	}
	g.svc = s
	_ = g.svc.Install()
	// store version, installed paths to db
	if err = g.db.New(ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
}

func (g *Grafana) Uninstall() error {
	var err error
	// install to global path
	// all files that are going to be installed
	var pkg *store.InstalledPackage
	pkg, err = g.db.Get(g.Name())
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
	if g.svc == nil {
		g.svc, _ = g.service()
	}
	_ = g.svc.Uninstall()
	dlPath := getDlPath(g.Name(), g.version.String())
	return os.RemoveAll(dlPath)
}

func (g *Grafana) Update(version update.Version) error {
	g.version = version
	if err := g.Uninstall(); err != nil {
		return err
	}
	if err := g.Download(); err != nil {
		return err
	}
	return g.Install()
}

func (g *Grafana) Run(args ...string) error {
	if g.svc == nil {
		g.svc, _ = g.service()
	}
	return g.svc.Start()
}

func (g *Grafana) Backup() error {
	if g.svc == nil {
		g.svc, _ = g.service()
	}
	return nil
}

func (g *Grafana) RunStop() error {
	return g.svc.Stop()
}

func (g *Grafana) Dependencies() []dep.Component {
	return nil
}

func (g *Grafana) IsDev() bool {
	return true
}

func (g *Grafana) RepoUrl() update.RepositoryURL {
	return grafanaBaseRepo
}
