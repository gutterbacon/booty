package downloader

import (
	"io/ioutil"
)

const (
	mbyte = 1 << 20
	gbyte = 1 << 30
)

func IsEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
