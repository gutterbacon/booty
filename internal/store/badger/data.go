package badger

import (
	bhold "github.com/timshannon/badgerhold/v2"
	"go.amplifyedge.org/booty-v2/internal/store"

	"go.amplifyedge.org/booty-v2/internal/logging"
)

type DB struct {
	store *bhold.Store
}

func NewDB(logger logging.Logger, dir string) *DB {
	options := bhold.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	options.Truncate = true
	options.Logger = logger

	st, err := bhold.Open(options)
	if err != nil {
		logger.Fatalf("cannot create database: %v", err)
	}

	return &DB{
		st,
	}
}

func (d *DB) New(i *store.InstalledPackage) error {
	return d.store.Upsert(i.Name, i)
}

func (d *DB) BulkNew(pkgs []*store.InstalledPackage) (err error) {
	tx := d.store.Badger().NewTransaction(true)
	for _, p := range pkgs {
		if err = d.store.TxUpsert(tx, p.Name, p); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) Get(packageName string) (*store.InstalledPackage, error) {
	var ip store.InstalledPackage
	if err := d.store.FindOne(&ip, bhold.Where("Name").Eq(packageName)); err != nil {
		return nil, err
	}
	return &ip, nil
}

func (d *DB) List() ([]*store.InstalledPackage, error) {
	var ips []store.InstalledPackage
	var result []*store.InstalledPackage
	if err := d.store.Find(&ips, nil); err != nil {
		return nil, err
	}
	for _, i := range ips {
		result = append(result, &i)
	}
	return result, nil
}

func (d *DB) Delete(packageName string) (err error) {
	tx := d.store.Badger().NewTransaction(true)
	if err = d.store.TxDelete(tx, packageName, &store.InstalledPackage{}); err != nil {
		return err
	}
	return tx.Commit()
}
