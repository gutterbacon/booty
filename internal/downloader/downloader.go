package downloader

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mholt/archiver/v3"
)

const (
	mbyte = 1 << 20
	gbyte = 1 << 30
)

func Download(dlUrl string, targetDir string) error {
	u, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	filename := filepath.Base(u.Path)
	dlDir := filepath.Dir(targetDir)
	destPath := filepath.Join(dlDir, filename)
	if err = downloadFile(dlUrl, destPath, filename); err != nil {
		return err
	}
	if err = extractDownloadedFile(destPath, filename, targetDir); err != nil {
		return err
	}
	return nil
}

func downloadFile(dlUrl, target, filename string) error {
	// download client
	client := grab.NewClient()
	req, _ := grab.NewRequest(target, dlUrl)

	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	go func() {
		for {
			select {
			case <-t.C:
				sz := humanize(resp.Size)
				fmt.Printf("%s %.2f / %.2f %s (%.2f%%)\n",
					filename,
					humanize(resp.BytesComplete()),
					sz,
					totalSz(sz),
					100*resp.Progress())

			case <-resp.Done:
				break
			}
		}
	}()

	// check for errors
	if err := resp.Err(); err != nil {
		return err
	}
	return nil
}

func humanize(i int64) float64 {
	sz := float64(i) / mbyte
	if sz >= 1000 {
		sz /= gbyte
	}
	return sz
}

func totalSz(f float64) string {
	if f < gbyte {
		return "MB"
	}
	return "GB"
}

func extractDownloadedFile(srcPath, filename, targetDir string) error {
	var err error
	fileExt := filepath.Ext(filename)
	switch fileExt {
	case ".gz", ".xz", ".zip", ".br", ".sz", ".bz2":
		err = archiver.Unarchive(srcPath, targetDir)
		if err != nil {
			return err
		}
		return os.RemoveAll(srcPath)
	default:
		return nil
	}
}

func GitClone(fetchUrl string, targetDir string, tag string) error {
	cloneOpts := &git.CloneOptions{
		URL:      fetchUrl,
		Progress: os.Stdout,
	}
	if tag != "" {
		cloneOpts.ReferenceName = plumbing.NewTagReferenceName(tag)
	}

	_, err := git.PlainClone(targetDir, false, cloneOpts)
	return err
}

func GitCheckout(tag string, targetDir string) error {
	r, err := git.PlainOpen(targetDir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	if err = w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		return err
	}
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
}
