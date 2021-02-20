package downloader

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mholt/archiver/v3"
)

func Download(dlUrl string, targetDir string) error {
	u, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	filename := filepath.Base(u.Path)
	dlDir := filepath.Dir(targetDir)
	destPath := filepath.Join(dlDir, filename)
	if err = downloadFile(dlUrl, destPath); err != nil {
		return err
	}
	if err = extractDownloadedFile(destPath, filename, targetDir); err != nil {
		return err
	}
	return nil
}

func downloadFile(dlUrl, target string) error {
	// download client
	client := grab.NewClient()
	req, _ := grab.NewRequest(target, dlUrl)

	fmt.Printf("Downloading %s\n", req.URL())
	resp := client.Do(req)
	fileSize := resp.Size

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	bar := pb.Full.Start64(fileSize)
	bar.SetMaxWidth(100)
	bar.Set(pb.Bytes, true)
	bar.Start()

	prevCompleted := int64(0)

	go func() {
		for {
			select {
			case <-t.C:
				completedNow := resp.BytesComplete()
				bar.Add64(completedNow - prevCompleted)
				prevCompleted = completedNow

			case <-resp.Done:
				bar.Finish()
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
