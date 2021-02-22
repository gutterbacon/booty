package downloader

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"io/ioutil"
	"os"
)

const (
	mbyte = 1 << 20
	gbyte = 1 << 30
)

func isEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
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
