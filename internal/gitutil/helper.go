package gitutil

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	remoteTemplate = `git@{{ .RepoHost }}:{{ .UpstreamOwner }}/{{ .Name }}`
)

type GitHelper struct {
	db        store.RepoStorer
	userEmail string
}

func NewHelper(db store.RepoStorer, userEmail string) *GitHelper {
	return &GitHelper{db: db, userEmail: userEmail}
}

func (gh *GitHelper) publicKey() (transport.AuthMethod, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	cb, err := ssh.NewSSHAgentAuth(currentUser.Name)
	return cb, err
}

func (gh *GitHelper) CatchupFork() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	err = r.Fetch(&git.FetchOptions{RemoteName: "upstream"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	err = osutil.Exec("git", "merge", "upstream/master")
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func (gh *GitHelper) StageAll() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	return wt.AddGlob(".")
}

func (gh *GitHelper) Stage(args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	for _, i := range args {
		_, err = wt.Add(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gh *GitHelper) Commit(msg string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Commit(msg, &git.CommitOptions{
		All:       false,
	})
	return err
}

func (gh *GitHelper) Push() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	r, err := gh.openGitDir(wd)
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{})
}

func (gh *GitHelper) SubmitPR() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	info, err := gh.RepoInfo(wd)
	if err != nil {
		return err
	}
	if info.Upstream == "" {
		return fmt.Errorf("empty upstream url")
	}
	var executable string
	switch osutil.GetOS() {
	case "darwin":
		executable = "open"
	case "linux":
		executable = "xdg-open"
	case "windows":
		executable = "start"
	}
	host := strings.Split(info.RepoHost, "-")[0]
	compareUrl := fmt.Sprintf("https://%s/%s/%s/compare/master..%s:%s", host, info.UpstreamOwner, info.Name, info.UserName, info.CurrentBranch)
	return osutil.Exec(executable, compareUrl)
}

func (gh *GitHelper) RegisterRepos(dirs ...string) error {
	// iterate over directories specified
	// register repos to the db
	for _, d := range dirs {
		abspath, err := filepath.Abs(d)
		if err != nil {
			return err
		}
		repoName := filepath.Base(d)
		if err = gh.db.RegisterRepo(repoName, abspath); err != nil {
			return err
		}
	}
	return nil
}

func (gh *GitHelper) SetupFork(upstreamOwner string) error {
	repos, err := gh.db.ListRepo()
	if err != nil {
		return err
	}
	for _, v := range repos {
		info, err := gh.RepoInfo(v)
		if err != nil {
			return err
		}
		r, err := gh.openGitDir(v)
		if err != nil {
			return err
		}
		if info.Upstream == "" {
			// setup upstream remotes
			info.UpstreamOwner = upstreamOwner
			t := template.Must(template.New("upstreamUrl").Parse(remoteTemplate))
			buf := new(bytes.Buffer)
			if err = t.ExecuteTemplate(buf, "upstreamUrl", info); err != nil {
				return err
			}
			info.Upstream = buf.String()
			_, err = r.CreateRemote(&config.RemoteConfig{
				Name: "upstream",
				URLs: []string{
					info.Upstream,
				},
			})
			if err != nil {
				return err
			}
		}
		// setup user
		cfg, err := r.Config()
		if err != nil {
			return err
		}
		if cfg.User.Name == "" {
			cfg.User.Name = info.UserName
		}
		if cfg.User.Email == "" && gh.userEmail != "" {
			cfg.User.Email = gh.userEmail
		}
		if err = r.SetConfig(cfg); err != nil {
			return err
		}
	}
	return nil
}

func (gh *GitHelper) CatchupAll() error {
	repos, err := gh.db.ListRepo()
	if err != nil {
		return err
	}
	for _, v := range repos {
		err = os.Chdir(v)
		if err != nil {
			return err
		}
		if err = gh.CatchupFork(); err != nil {
			return err
		}
	}
	return nil
}
