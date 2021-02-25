package file

import (
	"bytes"
	"encoding/gob"
	"go.amplifyedge.org/booty-v2/internal/errutil"
)

func marshalGob(repos map[string]string) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(&repos); err != nil {
		return nil, errutil.New(errutil.ErrOpenFile, err)
	}
	return b.Bytes(), nil
}

func (d *DB) readRepoInfoFile() (map[string]string, error) {
	var err error
	entries := map[string]string{}
	b := make([]byte, d.size)
	f, err := openFile(d.filepath, false)
	if err != nil {
		return nil, err
	}
	if _, err = f.Read(b); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	if err = gob.NewDecoder(buf).Decode(&entries); err != nil {
		if err.Error() == "unexpected EOF" {
			return entries, nil
		}
		return nil, err
	}
	return entries, nil
}

func (d *DB) RegisterRepo(name string, dirpath string) error {
	entries, err := d.readRepoInfoFile()
	if err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// check if repo exists, overwrite it if it is.
	entries[name] = dirpath
	b, err := marshalGob(entries)
	if err != nil {
		return err
	}
	if err = d.f.Truncate(0); err != nil {
		return err
	}
	sz, err := d.f.WriteAt(b, 0)
	if err != nil {
		return err
	}
	d.size = int64(sz)
	return nil
}

func (d *DB) UnregisterRepo(name string) error {
	entries, err := d.readRepoInfoFile()
	if err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := entries[name]
	if !ok {
		return nil
	}
	delete(entries, name)
	b, err := marshalGob(entries)
	if err != nil {
		return err
	}
	if err = d.f.Truncate(0); err != nil {
		return err
	}
	sz, err := d.f.WriteAt(b, 0)
	if err != nil {
		return err
	}
	d.size = int64(sz)
	return nil
}

func (d *DB) UnregisterAll() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	err := d.f.Truncate(0)
	if err != nil {
		return err
	}
	b, err := marshalGob(map[string]string{})
	if err != nil {
		return err
	}
	sz, err := d.f.WriteAt(b, 0)
	if err != nil {
		return err
	}
	d.size = int64(sz)
	return nil
}

func (d *DB) GetRepo(name string) string {
	entries, err := d.readRepoInfoFile()
	if err != nil {
		return ""
	}
	return entries[name]
}

func (d *DB) ListRepo() (map[string]string, error) {
	return d.readRepoInfoFile()
}
