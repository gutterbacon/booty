package osutil

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func GetOS() string {
	return runtime.GOOS
}

func GetAltOs() string {
	switch GetOS() {
	case "linux":
		return "linux"
	case "darwin":
		return "mac"
	case "windows":
		return "windows"
	default:
		return ""
	}
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

func getInstallPrefix() string {
	u, _ := user.Current()
	switch strings.ToLower(GetOS()) {
	case "windows":
		return filepath.Join(`C:\\ProgramData`, "booty")
	case "linux":
		return fmt.Sprintf("%s/.local/%s", u.HomeDir, "booty")
	case "darwin":
		return fmt.Sprintf("%s/Library/Application Support/%s", u.HomeDir, "booty")
	default:
		return "/usr/local"
	}
}

func SetupDirs() (err error) {
	prefix := getInstallPrefix()
	dirs := []string{"bin", "etc", "data", "downloads", "include"}
	for i := range dirs {
		dirPath := filepath.Join(prefix, dirs[i])
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

func GetBinDir() string {
	return filepath.Join(getInstallPrefix(), "bin")
}

func GetGoPath() string {
	return os.Getenv("GOPATH")
}

func GetEtcDir() string {
	return filepath.Join(getInstallPrefix(), "etc")
}

func GetDataDir() string {
	return filepath.Join(getInstallPrefix(), "data")
}

func GetDownloadDir() string {
	return filepath.Join(getInstallPrefix(), "downloads")
}

func GetIncludeDir() string {
	return filepath.Join(getInstallPrefix(), "include")
}

func CurUserChown(dir string) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		return err
	}
	c := exec.Command("chown", fmt.Sprintf("%s:%s", u.Name, g.Name), dir)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

// ExecSudo executes a command under "sudo".
func ExecSudo(cmd string, args ...string) error {
	scmd := exec.Command("sudo", append([]string{cmd}, args...)...)
	scmd.Stdin = os.Stdin
	scmd.Stderr = os.Stderr
	scmd.Stdout = os.Stdout

	return scmd.Run()
}

// Exec
func Exec(cmd string, args ...string) error {
	scmd := exec.Command(cmd, args...)
	scmd.Stdin = os.Stdin
	scmd.Stderr = os.Stderr
	scmd.Stdout = os.Stdout

	return scmd.Run()
}

// DetectPreq detect prequisite dependencies (golang and flutter)
// if it doesn't exists in the system, then return error
func DetectPreq() error {
	var executableName = "go"
	if GetOS() == "windows" {
		executableName += ".exe"
	}
	golangExists := ExeExists(executableName)
	gopathExists := checkEnv("GOPATH")
	if golangExists && gopathExists {
		return nil
	}
	return fmt.Errorf("golang executable doesn't exists, please install golang first")
}

func ExeExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkEnv(key string) bool {
	env := os.Getenv(key)
	return env != ""
}
