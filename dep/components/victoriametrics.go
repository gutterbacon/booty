package components

import (
	"embed"
	"fmt"
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

//go:emabed files/prometheus.yml
var prometheusCfgSample embed.FS

const (
	vicmetUrlFmt = "https://github.com/VictoriaMetrics/VictoriaMetrics"
)

type VicMet struct {
	version string
	db      *store.DB
	svcs    []*service.Svc
}

func NewVicMet(db *store.DB, version string) *VicMet {
	return &VicMet{
		version: version,
		db:      db,
	}
}

func (v *VicMet) service(promCfgPath, vmStoragePath string) ([]*service.Svc, error) {
	nameVer := fmt.Sprintf("%s-%s", v.Name(), v.version)
	vmConfig := &ks.Config{
		Name:        nameVer,
		DisplayName: nameVer,
		Description: "fast, cost-effective monitoring solution and time series database",
		Arguments: []string{
			"--promscape.config=" + promCfgPath,
			"--opentsdbListenAddr=:4242",
			"--httpListenAddr=:8428",
			"--storageDataPath=" + vmStoragePath,
		},
		Executable: filepath.Join(osutil.GetBinDir(), v.Name()),
		Option:     map[string]interface{}{},
	}
	vmAgentConfig := &ks.Config{
		Name:        "vmagent-" + v.version,
		DisplayName: "vmagent" + v.version,
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

func (v *VicMet) Version() string {
	return v.version
}

func (v *VicMet) Download() error {
	targetDir := getDlPath(v.Name(), v.version)
	return downloader.GitClone(vicmetUrlFmt, targetDir, "v"+v.version)
}

func (v *VicMet) Dependencies() []dep.Component {
	return nil
}

func (v *VicMet) Install() error {
	dlPath := getDlPath(v.Name(), v.version)
	var err error
	binDir := osutil.GetBinDir()
	if err = os.Chdir(dlPath); err != nil {
		return err
	}
	if err = osutil.Exec("make", "all"); err != nil {
		return err
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
		Version: v.version,
		FilesMap: map[string]int{
			filepath.Join(osutil.GetDataDir(), v.Name(), "storage"): 0755,
		},
	}
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, "bin", v.Name()):    {filepath.Join(binDir, v.Name()), 0755},
		filepath.Join(dlPath, "bin", "vmagent"):   {filepath.Join(binDir, "vmagent"), 0755},
		filepath.Join(dlPath, "bin", "vmalert"):   {filepath.Join(binDir, "vmalert"), 0755},
		filepath.Join(dlPath, "bin", "vmauth"):    {filepath.Join(binDir, "vmauth"), 0755},
		filepath.Join(dlPath, "bin", "vmbackup"):  {filepath.Join(binDir, "vmbackup"), 0755},
		filepath.Join(dlPath, "bin", "vmrestore"): {filepath.Join(binDir, "vmrestore"), 0755},
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
	return nil
}

func (v *VicMet) Run(args ...string) error {
	for _, s := range v.svcs {
		if err := s.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (v *VicMet) Update(version string) error {
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
