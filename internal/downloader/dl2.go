package downloader

import (
	"github.com/hashicorp/go-getter"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"net/url"
	"path/filepath"
	"sync"
)

func Download(dlUrl string, targetDir string) error {
	// make sure url is valid
	_, err := url.Parse(dlUrl)
	if err != nil {
		return err
	}
	if osutil.GetOS() == "windows" {
		return getter.Get(targetDir, dlUrl)
	}
	var notex bool
	if notex, err = isEmptyDir(targetDir); notex || err != nil {
		progBar := progressBar{}
		return getter.Get(targetDir, dlUrl, getter.WithProgress(&progBar))
	}
	return nil
}

//// defaultProgressBar is the default instance of a cheggaaa
//// progress bar.
//var defaultProgressBar getter.ProgressTracker = &progressBar{}

// ProgressBar wraps a github.com/cheggaaa/pb.Pool
// in order to display download progress for one or multiple
// downloads.
//
// If two different instance of ProgressBar try to
// display a progress only one will be displayed.
// It is therefore recommended to use DefaultProgressBar
type progressBar struct {
	// lock everything below
	lock sync.Mutex

	pool *pb.Pool

	pbs int
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
