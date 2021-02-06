package downloader

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/mholt/archiver/v3"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
)

func Download(dlUrl string, targetDir string) error {
	u, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	filename := filepath.Base(u.Path)
	destPath := filepath.Join(targetDir, filename)
	if err = downloadFile(dlUrl, targetDir); err != nil {
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

	fmt.Printf("Downloading %s", req.URL())
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
