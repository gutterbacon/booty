package osdetect

import (
	"runtime"
)

func GetOS() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

// We support only x86_64 or arm64 only for now
func GetAltArch() string {
	switch GetArch() {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "arm64v8"
	default:
		return ""
	}
}

