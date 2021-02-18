// Config contains version information
// of all the components
package config

import (
	"embed"
	"encoding/json"

	"go.amplifyedge.org/booty-v2/internal/logging"
)

//go:embed config.reference.json

var DefaultConfig embed.FS

type BinaryInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type AppConfig struct {
	DevMode  bool         `json:"dev"`
	Binaries []BinaryInfo `json:"binaries,omitempty"`
}

func NewAppConfig(logger logging.Logger, r []byte) *AppConfig {
	var ac AppConfig
	if err := json.Unmarshal(r, &ac); err != nil {
		logger.Fatalf("error parsing version information: %v", err)
	}
	if ac.Binaries == nil {
		ac.Binaries = []BinaryInfo{}
	}
	return &ac
}

func (ac *AppConfig) GetVersion(pkgName string) string {
	for _, pkg := range ac.Binaries {
		if pkgName == pkg.Name {
			return pkg.Version
		}
	}
	return ""
}
