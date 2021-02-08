package fileutil

import (
	"io"
	"os"
	"strings"

	"github.com/otiai10/copy"
)

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if sourceFileStat.Mode().IsRegular() {
		source, err := os.Open(src)
		if err != nil {
			return err
		}

		destination, err := os.Create(dst)
		if err != nil {
			return err
		}

		_, err = io.Copy(destination, source)
		_ = source.Close()
		_ = destination.Close()
		return err

	}

	opt := copy.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Shallow
		},
	}

	return copy.Copy(src, dst, opt)
}
