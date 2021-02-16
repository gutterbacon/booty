// update package provides a way to check for github release update
// it will parse the tag of the repository in question and compare it with the current version
// in the database, if there's update available, it will store it inside the badger database,
// and the agent will update the relevant component in question
package update

import (
	"context"
	"github.com/google/go-github/github"
	"time"

	"go.amplifyedge.org/booty-v2/internal/logging"
	"go.amplifyedge.org/booty-v2/internal/store"
)

const (
	defaultTimeout = 5 * time.Second
)

// RepositoryURL is the github repository url we will monitor
type RepositoryURL string

// Version is the version information for the releases.
type Version string

func (v Version) String() string {
	return string(v)
}

type Checker struct {
	logger     logging.Logger
	db         *store.DB                          // the badger database
	repos      map[RepositoryURL]Version          // repository urls
	updateFunc func(RepositoryURL, Version) error // call this function on new release
}

func NewChecker(logger logging.Logger, db *store.DB, repos map[RepositoryURL]Version, updateFunc func(r RepositoryURL, v Version) error) *Checker {
	return &Checker{
		db:         db,
		repos:      repos,
		logger:     logger,
		updateFunc: updateFunc,
	}
}

func (c *Checker) getRepoInfos() (chan *repoInfo, error) {
	rchan := make(chan *repoInfo, len(c.repos))
	for r, v := range c.repos {
		rinfo, err := parseGithubUrl(r, v)
		if err != nil {
			c.logger.Errorf("error parsing github url for %s: %v", r, err)
			return nil, err
		}
		rchan <- rinfo
	}
	close(rchan)
	return rchan, nil
}

func (c *Checker) CheckNewReleases() error {
	ghc := github.NewClient(nil)
	rinfos, err := c.getRepoInfos()
	if err != nil {
		return err
	}
	for k := range rinfos {
		if err = c.fetchLatest(k, ghc); err != nil {
			return err
		}
	}
	return nil
}

func (c *Checker) fetchLatest(r *repoInfo, ghc *github.Client) error {
	c.logger.Infof("checking update for: %s", r.repoUrl)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	release, _, err := ghc.Repositories.GetLatestRelease(ctx, r.repoUser, r.repoName)
	if err != nil {
		return err
	}
	v := *release.TagName
	if isTagNewer(r.curVer.toSemver(), parseReleaseTag(v)) {
		c.logger.Infof("latest version is newer: %s than current: %s, updating...", r.curVer, Version(v))
		// assign it to the repos
		c.repos[r.repoUrl] = Version(v)
		// do the update function
		if err = c.updateFunc(r.repoUrl, Version(v)); err != nil {
			return err
		}
	}
	return nil
}

func GetLatestVersion(repoUrl RepositoryURL) (string, error) {
	ghc := github.NewClient(nil)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	info, err := parseGithubUrl(repoUrl, "")
	if err != nil {
		return "", err
	}
	release, _, err := ghc.Repositories.GetLatestRelease(ctx, info.repoUser, info.repoName)
	if err != nil {
		return "", err
	}
	return *release.TagName, nil
}
