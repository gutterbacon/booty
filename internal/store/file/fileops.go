package file

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.amplifyedge.org/booty-v2/internal/errutil"
	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/store"
)

type allInstalledPackages struct {
	Packages []*store.InstalledPackage `json:"packages"`
}

var mu sync.Mutex

// It's a file db
type DB struct {
	mu          *sync.Mutex // mutex
	filepath    string      // path to file
	initialSize int64
	f           *os.File
	size        int64
	logger      logging.Logger // logger
}

func (d *DB) New(ip *store.InstalledPackage) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	var sz int
	allPkgs, err := d.getAllPkgs()
	if err != nil {
		return err
	}
	// truncate file to empty
	if err = d.f.Truncate(0); err != nil {
		return err
	}
	// check if it exists
	exists := false
	for _, p := range allPkgs.Packages {
		if p.Name == ip.Name {
			exists = true
			p = ip
		}
	}
	if !exists {
		allPkgs.Packages = append(allPkgs.Packages, ip)
	}

	// replace the file content
	var b []byte
	b, err = json.Marshal(&allPkgs)
	if err != nil {
		return err
	}
	sz, err = d.f.WriteAt(b, 0)
	if err != nil {
		return err
	}
	d.size = int64(sz)
	return nil
}

func (d *DB) BulkNew(packages []*store.InstalledPackage) error {
	for _, p := range packages {
		if err := d.New(p); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) Get(pkgName string) (*store.InstalledPackage, error) {
	if d.size == 0 {
		return nil, errutil.New(errutil.ErrEmptyFile, fmt.Errorf("no package installed of name: %s", pkgName))
	}
	allPkgs, err := d.getAllPkgs()
	if err != nil {
		return nil, err
	}
	for _, p := range allPkgs.Packages {
		if p.Name == pkgName {
			return p, nil
		}
	}
	return nil, fmt.Errorf("error: package not found")
}

func (d *DB) List() ([]*store.InstalledPackage, error) {
	allPkgs, err := d.getAllPkgs()
	if err != nil {
		return nil, err
	}
	return allPkgs.Packages, nil
}

func (d *DB) Delete(pkgName string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	allPkgs, err := d.getAllPkgs()
	if err != nil {
		return err
	}
	var pkgs []*store.InstalledPackage
	for _, p := range allPkgs.Packages {
		if p.Name == pkgName {
			continue
		}
		pkgs = append(pkgs, p)
	}
	allPkgs.Packages = pkgs
	err = d.f.Truncate(0)
	if err != nil {
		return err
	}
	b, err := json.Marshal(&allPkgs)
	sz, err := d.f.WriteAt(b, 0)
	if err != nil {
		return err
	}
	d.size = int64(sz)
	return nil
}

func (d *DB) getAllPkgs() (*allInstalledPackages, error) {
	byteData := make([]byte, d.size)
	f, err := openFile(d.filepath, false)
	if err != nil {
		return nil, err
	}
	if _, err = f.Read(byteData); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	var allPkgs allInstalledPackages
	if err = json.Unmarshal(byteData, &allPkgs); err != nil {
		return nil, err
	}
	return &allPkgs, nil
}

func NewDB(logger logging.Logger, fpath string, repoMode bool) (*DB, error) {
	f, size, err := newOrExistingWrite(&mu, fpath, repoMode)
	if err != nil {
		return nil, err
	}
	return &DB{
		mu:          &mu,
		filepath:    fpath,
		f:           f,
		size:        size,
		initialSize: size,
		logger:      logger,
	}, nil
}

// open new or existing file
func newOrExistingWrite(mu *sync.Mutex, fpath string, repoMode bool) (*os.File, int64, error) {
	mu.Lock()
	defer mu.Unlock()
	f, err := openFile(fpath, true)
	if err != nil {
		return nil, 0, errutil.New(errutil.ErrOpenFile, err)
	}
	info, err := f.Stat()
	if err != nil {
		return nil, 0, errutil.New(errutil.ErrOpenFile, err)
	}
	size := info.Size()
	if size == 0 && !repoMode {
		allPkgs := &allInstalledPackages{Packages: []*store.InstalledPackage{}}
		b, err := json.Marshal(&allPkgs)
		if err != nil {
			return nil, 0, errutil.New(errutil.ErrOpenFile, err)
		}
		wlen, err := f.WriteAt(b, 0)
		if err != nil {
			return nil, 0, errutil.New(errutil.ErrOpenFile, err)
		}
		size = int64(wlen)
	} else if size == 0 && repoMode {
		allRepos := map[string]string{}
		b, err := marshalGob(allRepos)
		wlen, err := f.WriteAt(b, 0)
		if err != nil {
			return nil, 0, errutil.New(errutil.ErrOpenFile, err)
		}
		size = int64(wlen)
	}
	return f, size, err
}

func openFile(fpath string, write bool) (*os.File, error) {
	mode := os.FileMode(0600)
	var flag int
	if write {
		flag = os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_CREATE | os.O_RDONLY
	}
	return os.OpenFile(fpath, flag, mode)
}
