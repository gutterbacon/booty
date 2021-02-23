package gitutil

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Name          string
	Directory     string
	Origin        string
	Upstream      string
	LastTag       string
	CurrentBranch string
	CurrentRef    string
	UserName      string
	UserEmail     string
	RepoHost      string
	OriginOwner   string
	UpstreamOwner string
}

func (r *RepoInfo) String() string {
	return fmt.Sprintf(`
	Repository Name:    %s
	Repository Dir :    %s
	# Upstream Info
	Upstream URL   :    %s
	Upstream Owner :    %s
	# Origin Info
	Origin URL     :    %s
	Origin Owner   :    %s
	# Current git directory information
	User           :    %s
	Email          :    %s
	Host           :    %s
	# Tags, Refs, Branches
	Last Tag       :    %s
	Current Ref    :    %s
	Current Branch :    %s
`, r.Name, r.Directory, r.Upstream, r.UpstreamOwner,
		r.Origin, r.OriginOwner, r.UserName, r.UserEmail, r.RepoHost,
		r.LastTag, r.CurrentRef, r.CurrentBranch,
	)
}

func (gh *GitHelper) openGitDir(dirpath string) (*git.Repository, error) {
	root, err := gitRoot(dirpath)
	if err != nil {
		return nil, err
	}
	return git.PlainOpen(root)
}

func (gh *GitHelper) RepoInfo(dirpath string) (*RepoInfo, error) {
	r, err := gh.openGitDir(dirpath)
	if err != nil {
		return nil, err
	}
	origin, err := r.Remote("origin")
	if err != nil {
		return nil, err
	}
	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}
	var repoName string
	var repoHost string
	var upstreamUrl string
	var lastTag string
	var originOwner string
	var upstreamOwner string
	userName := cfg.User.Name
	userEmail := cfg.User.Email
	originUrl := origin.Config().URLs[0]
	ep, err := transport.NewEndpoint(originUrl)
	if err != nil {
		return nil, err
	}
	if userName == "" && originUrl != "" {
		repo := strings.Split(ep.Path, "/")
		if len(repo) == 2 {
			userName = repo[0]
			repoName = repo[1]
			repoName = strings.TrimSuffix(repoName, ".git")
		}
	}
	originOwner = userName
	repoHost = ep.Host
	upstream, err := r.Remote("upstream")
	if err != nil {
		upstreamUrl = ""
	}
	if upstream != nil {
		upstreamUrl = upstream.Config().URLs[0]
	}
	if upstreamUrl != "" {
		ep, err = transport.NewEndpoint(upstreamUrl)
		if err != nil {
			return nil, err
		}
		repo := strings.Split(ep.Path, "/")
		if len(repo) == 2 {
			upstreamOwner = repo[0]
		}
	}
	tags, err := r.Tags()
	if err != nil {
		lastTag = ""
	}
	tagDoneChan := make(chan bool, 1)
	go func() {
		for {
			ref, err := tags.Next()
			if err != nil {
				tagDoneChan <- true
				return
			}
			lastTag = ref.Name().Short()
		}
	}()
	if <-tagDoneChan {
	}
	headRef, err := r.Head()
	if err != nil {
		return nil, err
	}
	currentRef := headRef.Hash().String()
	currentBranch := headRef.Name().Short()
	return &RepoInfo{
		Name:          repoName,
		Directory:     dirpath,
		Origin:        originUrl,
		Upstream:      upstreamUrl,
		LastTag:       lastTag,
		CurrentBranch: currentBranch,
		CurrentRef:    currentRef,
		UserName:      userName,
		UserEmail:     userEmail,
		RepoHost:      repoHost,
		OriginOwner:   originOwner,
		UpstreamOwner: upstreamOwner,
	}, nil
}

func gitRoot(path string) (string, error) {
	// normalize the path
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		fi, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			if !fi.IsDir() {
				return "", fmt.Errorf(".git exist but is not a directory")
			}
			return filepath.Join(path, ".git"), nil
		}
		if !os.IsNotExist(err) {
			// unknown error
			return "", err
		}

		// detect bare repo
		ok, err := isGitDir(path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}

		if parent := filepath.Dir(path); parent == path {
			return "", fmt.Errorf(".git not found")
		} else {
			path = parent
		}
	}
}

func isGitDir(path string) (bool, error) {
	markers := []string{"HEAD", "objects", "refs"}

	for _, marker := range markers {
		_, err := os.Stat(filepath.Join(path, marker))
		if err == nil {
			continue
		}
		if !os.IsNotExist(err) {
			// unknown error
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}
