package osutil

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
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

func GetInstallPrefix() string {
	switch strings.ToLower(GetOS()) {
	case "windows":
		return `C:\\ProgramData`
	case "linux":
		u, _ := user.Current()
		return fmt.Sprintf("%s/.local", u.HomeDir)
	case "darwin":
		return "/usr/local"
	default:
		return "/usr/local"
	}
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