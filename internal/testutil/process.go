package testutil

import (
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"golang.org/x/sys/unix"
)

func Kill() error {
	if osutil.GetOS() == "linux" || osutil.GetOS() == "darwin" {
		return unix.Kill(unix.Getpid(), unix.SIGTERM)
	}
	return nil
}
