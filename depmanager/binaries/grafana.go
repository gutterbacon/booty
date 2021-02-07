package binaries

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/pkg/helpers/fileutil"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/pkg/helpers/downloader"
	"go.amplifyedge.org/booty-v2/pkg/helpers/osdetect"
)

const (
	// version -- os-arch
	fetchUrlFormat = "https://dl.grafana.com/oss/release/grafana-%s.%s.tar.gz"
)

// Grafana implements both Versioner interface and Component interface
type Grafana struct {
	version string
	dlPath  string
}

func NewGrafana(version string) *Grafana {
	return &Grafana{version, ""}
}

// Gets grafana's version
func (g *Grafana) Version() string {
	return g.version
}

func (g *Grafana) Name() string {
	return "grafana"
}

func (g *Grafana) Download(targetDir string) error {
	// create target dir
	_ = os.MkdirAll(targetDir, 0755)

	osname := fmt.Sprintf("%s-%s", strings.ToLower(osdetect.GetOS()), osdetect.GetArch())
	fetchUrl := fmt.Sprintf(fetchUrlFormat, g.version, osname)
	err := downloader.Download(fetchUrl, targetDir)
	if err != nil {
		return err
	}
	g.dlPath = filepath.Join(targetDir, "grafana-"+g.version)
	return nil
}

func (g *Grafana) Install() error {
	var err error
	// install to global path
	prefixDir := osdetect.GetInstallPrefix()
	switch strings.ToLower(osdetect.GetOS()) {
	case "linux", "darwin":
		// create bin directory under $PREFIX
		binDir := filepath.Join(prefixDir, "bin")
		_ = os.MkdirAll(binDir, os.ModeDir)
		if err = osdetect.CurUserChown(binDir); err != nil {
			return err
		}
		// create etc dir under $PREFIX
		etcDir := filepath.Join(prefixDir, "etc")
		_ = os.MkdirAll(etcDir, os.ModeDir)
		if err = osdetect.CurUserChown(etcDir); err != nil {
			return err
		}

		// all files that are going to be installed
		filesMap := map[string][]interface{}{
			filepath.Join(g.dlPath, "bin", "grafana-server"): {filepath.Join(binDir, "grafana-server"), 0755},
			filepath.Join(g.dlPath, "bin", "grafana-cli"):    {filepath.Join(binDir, "grafana-cli"), 0755},
			filepath.Join(g.dlPath, "conf", "defaults.ini"):  {filepath.Join(etcDir, "grafana.ini"), 0644},
			filepath.Join(g.dlPath, "conf", "sample.ini"):    {filepath.Join(etcDir, "grafana.sample.ini"), 0644},
		}

		// copy file to the bin directory
		for k, v := range filesMap {
			if err = fileutil.Copy(k, v[0].(string)); err != nil {
				return err
			}
			if err = os.Chmod(v[0].(string), v[1].(os.FileMode)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Grafana) Uninstall() error {
	var err error
	// install to global path
	prefixDir := osdetect.GetInstallPrefix()
	switch strings.ToLower(osdetect.GetOS()) {
	case "linux", "darwin":
		// create bin directory under $PREFIX
		binDir := filepath.Join(prefixDir, "bin")
		_ = os.MkdirAll(binDir, os.ModeDir)

		// create etc dir under $PREFIX
		etcDir := filepath.Join(prefixDir, "etc")
		_ = os.MkdirAll(etcDir, os.ModeDir)

		// all files that are going to be installed
		filesList := []string{
			filepath.Join(binDir, "grafana-server"),
			filepath.Join(binDir, "grafana-cli"),
			filepath.Join(etcDir, "grafana.ini"),
			filepath.Join(etcDir, "grafana.sample.ini"),
		}

		// copy file to the bin directory
		for _, file := range filesList {
			if err = osdetect.ExecSudo("rm", "-rf", file); err != nil {
				return err
			}
		}
	}
	// remove downloaded files
	return os.RemoveAll(g.dlPath)
}

func (g *Grafana) Run() error {
	return nil
}

func (g *Grafana) Backup() error {
	return nil
}
