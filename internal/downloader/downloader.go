package downloader

import (
	"github.com/go-git/go-git/v5/plumbing"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-getter"

	"go.amplifyedge.org/booty-v2/internal/osutil"
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
	return getter.Get(targetDir, dlUrl, getter.WithProgress(defaultProgressBar))
}

// defaultProgressBar is the default instance of a cheggaaa
// progress bar.
var defaultProgressBar getter.ProgressTracker = &progressBar{}

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
		cpb.pool.Start()
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
				cpb.pool.Stop()
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
