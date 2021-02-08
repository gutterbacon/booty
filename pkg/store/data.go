package store

import (
	bhold "github.com/timshannon/badgerhold/v2"
	"go.amplifyedge.org/booty-v2/pkg/logging"
)

type DB struct {
	store *bhold.Store
}

func NewDB(logger logging.Logger, dir string) *DB {
	options := bhold.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir

	store, err := bhold.Open(options)
	if err != nil {
		logger.Fatalf("cannot create database: %v", err)
	}

	return &DB{
		store,
	}
}

type InstalledPackage struct {
	Name     string         `badgerholdIndex:"key"` // name of the package
	Version  string         // package version
	FilesMap map[string]int // files installed
}

func (d *DB) New(i *InstalledPackage) error {
	return d.store.Upsert(i.Name, i)
}

func (d *DB) BulkNew(pkgs []*InstalledPackage) (err error) {
	tx := d.store.Badger().NewTransaction(true)
	for _, p := range pkgs {
		if err = d.store.TxUpsert(tx, p.Name, p); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) Get(packageName string) (*InstalledPackage, error) {
	var ip InstalledPackage
	if err := d.store.FindOne(&ip, bhold.Where("Name").Eq(packageName)); err != nil {
		return nil, err
	}
	return &ip, nil
}

func (d *DB) List(query *bhold.Query) ([]*InstalledPackage, error) {
	var ips []InstalledPackage
	var result []*InstalledPackage
	if err := d.store.Find(&ips, query); err != nil {
		return nil, err
	}
	for _, i := range ips {
		result = append(result, &i)
	}
	return result, nil
}

func (d *DB) Delete(query *bhold.Query, packageNames []string) (err error) {
	tx := d.store.Badger().NewTransaction(true)
	for _, name := range packageNames {
		if err = d.store.TxDelete(tx, name, &InstalledPackage{}); err != nil {
			return err
		}
	}
	return tx.Commit()
}
