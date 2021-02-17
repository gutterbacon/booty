package store

type Storer interface {
	New(*InstalledPackage) error
	BulkNew([]*InstalledPackage) error
	Get(string) (*InstalledPackage, error)
	List() ([]*InstalledPackage, error)
	Delete(string) error
}

type InstalledPackage struct {
	Name     string         `badgerholdIndex:"key" json:"name"` // name of the package
	Version  string         `json:"version"`                    // package version
	FilesMap map[string]int `json:"filesMap"`                   // files installed
}


