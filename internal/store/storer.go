package store

type Storer interface {
	New(*InstalledPackage) error
	BulkNew([]*InstalledPackage) error
	Get(string) (*InstalledPackage, error)
	List() ([]*InstalledPackage, error)
	Delete(string) error
}

type InstalledPackage struct {
	Name     string            `json:"name"` // name of the package
	Version  string            `json:"version"`                    // package version
	FilesMap map[string]string `json:"filesMap"`                   // files installed
}
