package components

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
	"os"
	"path/filepath"
)

const (
	jbRepo = "https://github.com/jsonnet-bundler/jsonnet-bundler"
)

type Jb struct {
	version update.Version
	db      *store.DB
}

func (j *Jb) IsDev() bool {
	return true
}

func NewJb(db *store.DB) *Jb {
	return &Jb{
		db: db,
	}
}

func (j *Jb) SetVersion(v update.Version) {
	j.version = v
}

func (j *Jb) Name() string {
	return "jb"
}

func (j *Jb) Version() update.Version {
	return j.version
}

func (j *Jb) Download() error {
	targetDir := getDlPath(j.Name(), j.version.String())
	if osutil.DirExists(targetDir) {
		return downloader.GitCheckout("v"+j.version.String(), targetDir)
	}
	return downloader.GitClone(jbRepo, targetDir, "")
}

func (j *Jb) Dependencies() []dep.Component {
	return nil
}

func (j *Jb) Install() error {
	targetDir := getDlPath(j.Name(), j.version.String())
	var err error
	binDir := osutil.GetBinDir()
	if err = os.Chdir(targetDir); err != nil {
		return err
	}
	if err = osutil.Exec("make", "static"); err != nil {
		return err
	}
	ip := store.InstalledPackage{
		Name:     j.Name(),
		Version:  j.version.String(),
		FilesMap: map[string]int{},
	}
	filesMap := map[string][]interface{}{
		filepath.Join(targetDir, "_output", j.Name()): {filepath.Join(binDir, j.Name()), 0755},
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
	if err = j.db.New(&ip); err != nil {
		return err
	}
	return nil
}

func (j *Jb) Uninstall() error {
	var err error
	var pkg *store.InstalledPackage
	pkg, err = j.db.Get(j.Name())
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
	return os.RemoveAll(getDlPath(j.Name(), j.version.String()))
}

func (j *Jb) Run(args ...string) error {
	return osutil.Exec(j.Name(), args...)
}

func (j *Jb) Update(version update.Version) error {
	j.SetVersion(version)
	if err := j.Uninstall(); err != nil {
		return err
	}
	if err := j.Download(); err != nil {
		return err
	}
	return j.Install()
}

func (j *Jb) RunStop() error {
	return nil
}

func (j *Jb) Backup() error {
	return nil
}

func (j *Jb) RepoUrl() update.RepositoryURL {
	return jbRepo
}
