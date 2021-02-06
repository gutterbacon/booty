package binaries

import (
	"fmt"
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
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}
	osname := fmt.Sprintf("%s-%s", strings.ToLower(osdetect.GetOS()), osdetect.GetArch())
	fetchUrl := fmt.Sprintf(fetchUrlFormat, g.version, osname)
	err = downloader.Download(fetchUrl, targetDir)
	if err != nil {
		return err
	}
	g.dlPath = filepath.Join(targetDir, "grafana-"+g.version)
	return nil
}

func (g *Grafana) Install() error {
	return nil
}

func (g *Grafana) Uninstall() error {
	return nil
}

func (g *Grafana) Run() error {
	return
}

func (g *Grafana) Backup() error {
	return nil
}
