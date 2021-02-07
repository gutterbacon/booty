package downloader

import (
	"net/url"

	"github.com/hashicorp/go-getter"
)

func Download(dlUrl string, targetDir string) error {
	// make sure url is valid
	_, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	return getter.Get(targetDir, dlUrl)
}
