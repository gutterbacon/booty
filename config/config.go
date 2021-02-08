// Config contains version information
// of all the components
package config

import (
	"encoding/json"
	"io/ioutil"

	"go.amplifyedge.org/booty-v2/pkg/logging"
)

type BinaryInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type VersionInfo struct {
	Binaries []BinaryInfo `json:"binaries"`
}

func NewVersionInfo(logger logging.Logger, jsonFile string) *VersionInfo {
	var vi VersionInfo
	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logger.Fatalf("error getting version information: %v", err)
	}
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
