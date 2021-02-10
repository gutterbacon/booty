package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.amplifyedge.org/booty-v2/dep"

	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
)

const (
	// version -- version -- os-arch except windows
	protocUrlFormat = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-%s.zip"
)

type Protoc struct {
	version      string
	db           *store.DB
	dependencies []dep.Component
}

func NewProtoc(db *store.DB, version string, deps []dep.Component) *Protoc {
	return &Protoc{
		version:      version,
		db:           db,
		dependencies: deps,
	}
}

func (p *Protoc) Name() string {
	return "protoc"
}

func (p *Protoc) Version() string {
	return p.version
}

func (p *Protoc) Download() error {
	if osutil.GetArch() != "amd64" {
		return fmt.Errorf("error: unsupported arch: %v", osutil.GetArch())
	}
	// download all dependencies
	errChan := make(chan error, len(p.dependencies))
	var wg sync.WaitGroup
	for i := 0; i < len(p.dependencies); i++ {
		wg.Add(1)
		j := i
		w := newWorkerType("download", osutil.GetDownloadDir(), p.dependencies, errChan)
		go func() {
			defer wg.Done()
			w.do(j)
		}()
	}
	wg.Wait()
	if err := <-errChan; err != nil {
		return err
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
	fmt.Printf("Downloading: %s Version: %s From: %s", p.Name(), p.Version(), fetchUrl)
	targetDir := getDlPath("protobuf", p.version)
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
	dlPath := getDlPath("protobuf", p.version)
	filesMap := map[string][]interface{}{
		filepath.Join(dlPath, "bin", executableName): {filepath.Join(binDir, executableName), 0755},
		filepath.Join(dlPath, "include", "google"):   {filepath.Join(includeDir, "google"), 0755},
	}

	ip := store.InstalledPackage{
		Name:     p.Name(),
		Version:  p.version,
		FilesMap: map[string]int{},
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
	if err = p.db.New(&ip); err != nil {
		return err
	}
	return os.RemoveAll(dlPath)
}

func (p *Protoc) Update(version string) error {
	p.version = version
	if err := p.Uninstall(); err != nil {
		return err
	}
	if err := p.Download(); err != nil {
		return err
	}
	return p.Install()
}

func (p *Protoc) Uninstall() error {
	var err error

	// uninstall all dependencies
	errChan := make(chan error, len(p.dependencies))
	var wg sync.WaitGroup
	for i := 0; i < len(p.dependencies); i++ {
		wg.Add(1)
		j := i
		w := newWorkerType("uninstall", "", p.dependencies, errChan)
		go func() {
			defer wg.Done()
			w.do(j)
		}()
	}
	wg.Wait()
	if err := <-errChan; err != nil {
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
