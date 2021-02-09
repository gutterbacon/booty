// Config contains version information
// of all the components
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"go.amplifyedge.org/booty-v2/pkg/logging"
)

type BinaryInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type VersionInfo struct {
	DevMode  bool         `json:"dev"`
	Binaries []BinaryInfo `json:"binaries"`
}

func NewVersionInfo(logger logging.Logger, r io.Reader) *VersionInfo {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Fatalf("error reading version information: %v", err)
	}
	var vi VersionInfo
	if err = json.Unmarshal(b, &vi); err != nil {
		logger.Fatalf("error parsing version information: %v", err)
	}
	return &vi
}

func (vi *VersionInfo) GetVersion(pkgName string) string {
	for _, pkg := range vi.Binaries {
		if pkgName == pkg.Name {
			return pkg.Version
		}
	}
	return ""
}
