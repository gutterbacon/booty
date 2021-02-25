package gitutil

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
)

func (gh *GitHelper) CreateTag(tagName string, tagMsg string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	ok, err := setTag(r, tagName, tagMsg)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("tag exists: %s", tagName)
	}
	return nil
}

func (gh *GitHelper) DeleteTag(tagName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	return r.DeleteTag(tagName)
}

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	tags, err := r.TagObjects()
	if err != nil {
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		return false
	}
	return res
}

func setTag(r *git.Repository, tag, msg string) (bool, error) {
	if tagExists(tag, r) {
		return false, nil
	}
	h, err := r.Head()
	if err != nil {
		return false, err
	}
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: msg,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (gh *GitHelper) pushTags(r *git.Repository, upstreamRemoteName string) error {

	auth, err := gh.publicKey()
	if err != nil {
		return err
	}

	po := &git.PushOptions{
		RemoteName: upstreamRemoteName,
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       auth,
	}
	err = r.Push(po)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil
		}
		return err
	}

	return nil
}

func (gh *GitHelper) PushTag() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	return gh.pushTags(r, "upstream")
}
