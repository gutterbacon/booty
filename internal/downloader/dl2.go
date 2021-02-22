package downloader

import (
	pb "github.com/cheggaaa/pb"
	"github.com/hashicorp/go-getter"
	"io"
	"net/url"
	"path/filepath"

	//"path/filepath"
	"sync"
)

func Download(dlUrl string, targetDir string) error {
	// make sure url is valid
	_, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	var notex bool

	httpGetter := getter.HttpGetter{Netrc: true}

	if notex, err = isEmptyDir(targetDir); notex || err != nil {
		pbar := progressBar{}
		progress := getter.WithProgress(&pbar)
		client := &getter.Client{
			Src:     dlUrl,
			Dst:     targetDir,
			Mode:    getter.ClientModeAny,
			Options: []getter.ClientOption{progress},
			Getters: map[string]getter.Getter{
				"file":  &getter.FileGetter{Copy: false},
				"http":  &httpGetter,
				"https": &httpGetter,
				"s3":    new(getter.S3Getter),
				"gcs":   new(getter.GCSGetter),
			},
		}
		return client.Get()
	}
	return nil
}

type progressBar struct {
	// lock everything below
	lock sync.Mutex
	pool *pb.Pool
	pbs  int
}

func ProgressBarConfig(bar *pb.ProgressBar, prefix string) {
	bar.SetUnits(pb.U_BYTES)
	bar.Prefix(prefix)
}

// TrackProgress instantiates a new progress bar that will
// display the progress of stream until closed.
// total can be 0.
func (cpb *progressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	newPb := pb.New64(totalSize)
	newPb.Set64(currentSize)
	ProgressBarConfig(newPb, filepath.Base(src))
	if cpb.pool == nil {
		cpb.pool = pb.NewPool()
		_ = cpb.pool.Start()
	}
	cpb.pool.Add(newPb)
	reader := newPb.NewProxyReader(stream)

	cpb.pbs++
	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			newPb.Finish()
			cpb.pbs--
			if cpb.pbs <= 0 {
				_ = cpb.pool.Stop()
				cpb.pool = nil
			}
			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
