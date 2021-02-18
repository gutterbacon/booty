package components

import (
	"embed"
	"fmt"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"io/ioutil"
	"os"
	"path/filepath"

	ks "github.com/kardianos/service"
	"go.amplifyedge.org/booty-v2/dep"

	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/service"
)

//go:embed files/Caddyfile
var caddyFileSample embed.FS

const (
	// version -- version -- os_arch
	caddyRepo      = "https://github.com/caddyserver/caddy"
	caddyUrlFormat = caddyRepo + "/releases/download/v%s/caddy_%s_%s.%s"
)

type Caddy struct {
	version update.Version
	db      store.Storer
	svc     *service.Svc
}

func (c *Caddy) IsService() bool {
	return true
}

func NewCaddy(db store.Storer) *Caddy {
	return &Caddy{
		db: db,
	}
}

func (c *Caddy) SetVersion(v update.Version) {
	c.version = v
}

func (c *Caddy) Version() update.Version {
	return c.version
}

func (c *Caddy) Name() string {
	return "caddy"
}

func (c *Caddy) service() (*service.Svc, error) {
	config := &ks.Config{
		Name:        c.Name(),
		DisplayName: c.Name(),
		Description: "Extensible platform that uses TLS by default",
		Arguments: []string{
			"run",
			"--environ",
			"--config",
			filepath.Join(osutil.GetEtcDir(), "caddy", "Caddyfile"),
		},
		Executable: filepath.Join(osutil.GetBinDir(), c.Name()),
		Option:     map[string]interface{}{},
	}
	return service.NewService(config)
}

func (c *Caddy) Download() error {
	downloadDir := getDlPath(c.Name(), c.version.String())
	_ = os.MkdirAll(downloadDir, 0755)
	var osname string
	var ext string
	switch osutil.GetOS() {
	case "linux", "darwin":
		ext = "tar.gz"
		osname = fmt.Sprintf("%s_%s", osutil.GetAltOs(), osutil.GetArch())
	case "windows":
		ext = "zip"
		osname = fmt.Sprintf("%s_%s", osutil.GetOS(), osutil.GetArch())
	}
	fetchUrl := fmt.Sprintf(caddyUrlFormat, c.version, c.version, osname, ext)
	err := downloader.Download(fetchUrl, downloadDir)
	if err != nil {
		return err
	}
	return nil
}

func (c *Caddy) Install() error {
	var err error
	// install to global path
	binDir := osutil.GetBinDir()
	dlPath := getDlPath(c.Name(), c.version.String())

	// all files that are going to be installed
	executableName := c.Name()
	switch osutil.GetOS() {
	case "windows":
		executableName += ".exe"
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, executableName): {filepath.Join(binDir, executableName), 0755},
	}

	// copy file to the global bin directory
	caddyConfigPath := filepath.Join(osutil.GetEtcDir(), "caddy", "Caddyfile")
	_ = os.MkdirAll(filepath.Dir(caddyConfigPath), 0755)
	ip, err := commonInstall(c, filesMap)
	if err != nil {
		return err
	}

	// install default config, only if the config doesn't exists
	// TODO: prompt user?
	if exists := osutil.Exists(caddyConfigPath); !exists {
		caddyData, err := caddyFileSample.ReadFile("files/Caddyfile")
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(caddyConfigPath, caddyData, 0600); err != nil {
			return err
		}
	}
	// install service
	s, err := c.service()
	if err != nil {
		return err
	}
	c.svc = s
	_ = c.svc.Install()
	if err = c.db.New(ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
}

func (c *Caddy) Uninstall() error {
	var err error
	dlPath := getDlPath(c.Name(), c.version.String())
	// install to global path

	// all files that are going to be installed
	var pkg *store.InstalledPackage
	pkg, err = c.db.Get(c.Name())
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
	if c.svc == nil {
		c.svc, _ = c.service()
	}
	_ = c.svc.Stop()
	_ = c.svc.Uninstall()
	// remove downloaded files
	return os.RemoveAll(dlPath)
}

func (c *Caddy) Update(version update.Version) error {
	c.SetVersion(version)
	if err := c.Uninstall(); err != nil {
		return err
	}
	if err := c.Download(); err != nil {
		return err
	}
	return c.Install()
}

func (c *Caddy) Run(_ ...string) error {
	if c.svc == nil {
		c.svc, _ = c.service()
	}
	return c.svc.Start()
}

func (c *Caddy) Backup() error {
	// TODO
	return nil
}

func (c *Caddy) RunStop() error {
	if c.svc == nil {
		c.svc, _ = c.service()
	}
	return c.svc.Stop()
}

func (c *Caddy) Dependencies() []dep.Component {
	return nil
}

func (c *Caddy) IsDev() bool {
	return false
}

func (c *Caddy) RepoUrl() update.RepositoryURL {
	return caddyRepo
}
