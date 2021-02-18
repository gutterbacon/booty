package fileutil

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/otiai10/copy"
	"golang.org/x/mod/sumdb/dirhash"
)

// Copy copies a file or directory from source to destination
// returning directory hash or file hash
func Copy(src, dst string) (string, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if sourceFileStat.Mode().IsRegular() {
		var source *os.File
		source, err = os.Open(src)
		if err != nil {
			return "", err
		}

		destination, err := os.Create(dst)
		if err != nil {
			return "", err
		}

		wlen, err := io.Copy(destination, source)
		if err != nil {
			return "", err
		}
		if wlen == 0 {
			return "", errors.New("not copying")
		}

		_ = source.Close()
		_ = destination.Close()

		source, err = os.Open(src)
		if err != nil {
			return "", err
		}

		srcStat, err := source.Stat()
		if err != nil {
			return "", err
		}

		sourceContent := make([]byte, srcStat.Size())
		_, err = source.Read(sourceContent)
		if err != nil {
			return "", err
		}

		sumBytes := sha256.Sum256(sourceContent)
		sum := fmt.Sprintf("%x", sumBytes)

		_ = source.Close()

		return sum, err

	}

	opt := copy.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Shallow
		},
	}

	if err = copy.Copy(src, dst, opt); err != nil {
		return "", err
	}

	return dirhash.HashDir(src, "", dirhash.Hash1)
}
