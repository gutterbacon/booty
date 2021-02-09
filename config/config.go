// Config contains version information
// of all the components
package config

import (
	"encoding/json"

	"go.amplifyedge.org/booty-v2/internal/logging"
)

type BinaryInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type AppConfig struct {
	DevMode  bool         `json:"dev"`
	Binaries []BinaryInfo `json:"binaries"`
}

func NewAppConfig(logger logging.Logger, r []byte) *AppConfig {
	var ac AppConfig
	if err := json.Unmarshal(r, &ac); err != nil {
		logger.Fatalf("error parsing version information: %v", err)
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

var DefaultConfig = AppConfig{
	DevMode: true,
	Binaries: []BinaryInfo{
		{
			Name:    "grafana",
			Version: "7.4.0",
		},
		{
			Name:    "goreleaser",
			Version: "0.155.1",
		},
		{
			Name:    "caddy",
			Version: "2.3.0",
		},
		{
			Name:    "protoc",
			Version: "3.14.0",
		},
		{
			Name:    "protoc-gen-go",
			Version: "1.25.0",
		},
		{
			Name:    "protoc-gen-cobra",
			Version: "0.4.1",
		},
		{
			Name:    "protoc-gen-go-grpc",
			Version: "1.1.0",
		},
	},
}
