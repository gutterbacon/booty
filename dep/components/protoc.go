package components

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
	"strings"

	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/osutil"
)

const (
	// version -- version -- os-arch except windows
	protocRepo      = "https://github.com/protocolbuffers/protobuf"
	protocUrlFormat = protocRepo + "/releases/download/v%s/protoc-%s-%s.zip"
)

type Protoc struct {
	version      update.Version
	db           store.Storer
	dependencies []dep.Component
}

func (p *Protoc) IsService() bool {
	return false
}

func NewProtoc(db store.Storer, deps []dep.Component) *Protoc {
	return &Protoc{
		db:           db,
		dependencies: deps,
	}
}

func (p *Protoc) Name() string {
	return "protoc"
}

func (p *Protoc) Version() update.Version {
	return p.version
}

func (p *Protoc) SetVersion(v update.Version) {
	p.version = v
}

func (p *Protoc) Download() error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	// download all dependencies
	for _, d := range p.dependencies {
		if err := d.Download(); err != nil {
			return err
		}
	}

	var osName string
	switch osutil.GetOS() {
	case "linux":
		osName = fmt.Sprintf("%s-%s", osutil.GetOS(), osutil.GetAltArch())
	case "darwin":
		osName = fmt.Sprintf("%s-%s", "osx", osutil.GetAltArch())
	case "windows":
		osName = "win64"
	}
	fetchUrl := fmt.Sprintf(protocUrlFormat, p.version, p.version, osName)
	targetDir := getDlPath("protobuf", p.version.String())
	err := downloader.Download(fetchUrl, targetDir)
	if err != nil {
		return err
	}
	return nil
}

func (p *Protoc) Install() error {
	var err error
	// install to path
	binDir := osutil.GetBinDir()
	includeDir := osutil.GetIncludeDir()

	executableName := p.Name()
	if osutil.GetOS() == "windows" {
		executableName += ".exe"
	}

	// install all dependencies
	for _, d := range p.dependencies {
		if err = d.Install(); err != nil {
			return err
		}
	}

	// all files that are going to be installed
	dlPath := getDlPath("protobuf", p.version.String())
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, "bin", executableName): {filepath.Join(binDir, executableName), 0755},
		filepath.Join(dlPath, "include", "google"):   {filepath.Join(includeDir, "google"), 0755},
	}

	// copy file to the bin directory
	ip, err := commonInstall(p, filesMap)
	if err != nil {
		return err
	}

	if err = p.db.New(ip); err != nil {
		return err
	}
	return nil
}

func (p *Protoc) Update(version update.Version) error {
	return commonUpdate(p, version)
}

func (p *Protoc) Uninstall() error {
	var err error

	// uninstall all dependencies
	for _, d := range p.dependencies {
		if err = d.Uninstall(); err != nil {
			return err
		}
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

func (p *Protoc) Run(args ...string) error {
	pkg, err := p.db.Get(p.Name())
	if err != nil {
		return err
	}
	var executable string
	for k := range pkg.FilesMap {
		if strings.Contains(k, p.Name()) {
			executable = k
		}
	}
	var arguments []string
	arguments = append(arguments, "-I", osutil.GetIncludeDir())
	arguments = append(arguments, args...)
	return osutil.Exec(executable, arguments...)
}

func (p *Protoc) RunStop() error {
	// We do not need to implement this
	return nil
}

func (p *Protoc) Backup() error {
	// We do not need to implement this.
	return nil
}

func (p *Protoc) Dependencies() []dep.Component {
	return p.dependencies
}

func (p *Protoc) IsDev() bool {
	return true
}

func (p *Protoc) RepoUrl() update.RepositoryURL {
	return protocRepo
}
