package components

import (
	"embed"
	"go.amplifyedge.org/booty-v2/internal/update"
	"io/ioutil"
	"os"
	"path/filepath"

	ks "github.com/kardianos/service"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/service"
	"go.amplifyedge.org/booty-v2/internal/store"
)

//go:embed files/prometheus.yml
var prometheusCfgSample embed.FS

const (
	vicmetUrlFmt = "https://github.com/VictoriaMetrics/VictoriaMetrics"
)

type VicMet struct {
	version update.Version
	db      *store.DB
	svcs    []*service.Svc
}

func NewVicMet(db *store.DB) *VicMet {
	return &VicMet{
		db: db,
	}
}

func (v *VicMet) SetVersion(ver update.Version) {
	v.version = ver
}

func (v *VicMet) service(promCfgPath, vmStoragePath string) ([]*service.Svc, error) {
	vmConfig := &ks.Config{
		Name:        v.Name(),
		DisplayName: v.Name(),
		Description: "fast, cost-effective monitoring solution and time series database",
		Arguments: []string{
			"--promscrape.config=" + promCfgPath,
			"--opentsdbListenAddr=:4242",
			"--httpListenAddr=:8428",
			"--storageDataPath=" + vmStoragePath,
		},
		Executable: filepath.Join(osutil.GetBinDir(), v.Name()),
		Option:     map[string]interface{}{},
	}
	vmAgentConfig := &ks.Config{
		Name:        "vmagent",
		DisplayName: "vmagent",
		Description: "victoria metrics agent",
		Arguments: []string{
			"--remoteWrite.url=http://localhost:8428/api/v1/write",
		},
		Executable: filepath.Join(osutil.GetBinDir(), "vmagent"),
		Option:     map[string]interface{}{},
	}
	vmetsvc, err := service.NewService(vmConfig)
	if err != nil {
		return nil, err
	}
	vmagentSvc, err := service.NewService(vmAgentConfig)
	if err != nil {
		return nil, err
	}
	return []*service.Svc{
		vmetsvc, vmagentSvc,
	}, nil
}

func (v *VicMet) Name() string {
	return "victoria-metrics"
}

func (v *VicMet) Version() update.Version {
	return v.version
}

func (v *VicMet) Download() error {
	// victoriametrics doesn't support windows yet
	if osutil.GetOS() == "windows" {
		return nil
	}
	targetDir := getDlPath(v.Name(), v.version.String())
	if osutil.DirExists(targetDir) {
		return downloader.GitCheckout("v"+v.version.String(), targetDir)
	}
	return downloader.GitClone(vicmetUrlFmt, targetDir, "v"+v.version.String())
}

func (v *VicMet) Dependencies() []dep.Component {
	return nil
}

func (v *VicMet) Install() error {
	if osutil.GetOS() == "windows" {
		return nil
	}
	dlPath := getDlPath(v.Name(), v.version.String())
	var err error
	binDir := osutil.GetBinDir()
	if err = os.Chdir(dlPath); err != nil {
		return err
	}
	recipes := []string{"victoria-metrics", "vmagent", "vmalert", "vmauth", "vmbackup", "vmctl", "vminsert", "vmrestore", "vmselect", "vmstorage"}
	// make pure local binaries
	filesMap := map[string][]interface{}{}
	for _, r := range recipes {
		if err = osutil.Exec("make", "app-local-pure", "APP_NAME="+r); err != nil {
			return err
		}
		filesMap[filepath.Join(dlPath, "bin", r+"-pure")] = []interface{}{
			filepath.Join(binDir, r), 0755,
		}
	}
	vmStoragePath := filepath.Join(osutil.GetDataDir(), v.Name(), "storage")
	if err = os.MkdirAll(vmStoragePath, 0755); err != nil {
		return err
	}
	vmConfigPath := filepath.Join(osutil.GetEtcDir(), v.Name(), "prometheus.yml")
	if err = os.MkdirAll(vmConfigPath, 0755); err != nil {
		return err
	}
	ip := store.InstalledPackage{
		Name:    v.Name(),
		Version: v.version.String(),
		FilesMap: map[string]int{
			filepath.Join(osutil.GetDataDir(), v.Name(), "storage"): 0755,
		},
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
	// install default config
	if exists := osutil.Exists(vmConfigPath); !exists {
		promData, err := prometheusCfgSample.ReadFile("files/prometheus.yml")
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(vmConfigPath, promData, 0600); err != nil {
			return err
		}
	}
	// install services
	svcs, err := v.service(vmConfigPath, vmStoragePath)
	if err != nil {
		return err
	}
	v.svcs = svcs
	for _, s := range v.svcs {
		if err = s.Install(); err != nil {
			return err
		}
	}
	if err = v.db.New(&ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
}

func (v *VicMet) Uninstall() error {
	if osutil.GetOS() == "windows" {
		return nil
	}
	var err error
	var pkg *store.InstalledPackage
	pkg, err = v.db.Get(v.Name())
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
	for _, s := range v.svcs {
		if err = s.Uninstall(); err != nil {
			return err
		}
	}
	return nil
}

func (v *VicMet) Run(args ...string) error {
	if osutil.GetOS() == "windows" {
		return nil
	}
	vmStoragePath := filepath.Join(osutil.GetDataDir(), v.Name(), "storage")
	vmConfigPath := filepath.Join(osutil.GetEtcDir(), v.Name(), "prometheus.yml")
	if v.svcs == nil {
		svcs, err := v.service(vmConfigPath, vmStoragePath)
		if err != nil {
			return err
		}
		v.svcs = svcs
	}
	for _, s := range v.svcs {
		if err := s.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (v *VicMet) Update(version update.Version) error {
	if osutil.GetOS() == "windows" {
		return nil
	}
	v.version = version
	if err := v.Uninstall(); err != nil {
		return err
	}
	if err := v.Download(); err != nil {
		return err
	}
	return v.Install()
}

func (v *VicMet) RunStop() error {
	if osutil.GetOS() == "windows" {
		return nil
	}
	vmStoragePath := filepath.Join(osutil.GetDataDir(), v.Name(), "storage")
	vmConfigPath := filepath.Join(osutil.GetEtcDir(), v.Name(), "prometheus.yml")
	if v.svcs == nil {
		svcs, err := v.service(vmConfigPath, vmStoragePath)
		if err != nil {
			return err
		}
		v.svcs = svcs
	}
	for _, s := range v.svcs {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (v *VicMet) Backup() error {
	// TODO
	return nil
}

func (v *VicMet) IsDev() bool {
	return true
}

func (v *VicMet) RepoUrl() update.RepositoryURL {
	return vicmetUrlFmt
}
