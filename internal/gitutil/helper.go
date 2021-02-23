package gitutil

import (
	"bytes"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	remoteTemplate = `git@{{ .RepoHost }}:{{ .UpstreamOwner }}/{{ .RepoName }}`
)

type GitHelper struct {
	db         store.RepoStorer
	sshKeyPath string
}

func NewHelper(db store.RepoStorer, sshKeyPath string) *GitHelper {
	return &GitHelper{db: db, sshKeyPath: sshKeyPath}
}

func publicKey(filePath string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshKey, _ := ioutil.ReadFile(filePath)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
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
	if err != nil {
		return err
	}
	if err = osutil.Exec("git", "merge", "upstream/master"); err != nil {
		return err
	}
	return nil
}

func (gh *GitHelper) Stage() error {
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
	return wt.AddGlob("*")
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
	_, err = wt.Commit(msg, nil)
	return err
}

func (gh *GitHelper) SubmitPR(prMsg string) error {
	return nil
}

func (gh *GitHelper) CreateTag(tagName string, tagMsg string) error {
	return nil
}

func (gh *GitHelper) DeleteTag(tagName string) {
	return
}

func (gh *GitHelper) RegisterRepos(dirs ...string) error {
	// iterate over directories specified
	var err error
	// register repos to the db
	for _, d := range dirs {
		repoName := filepath.Base(d)
		if err = gh.db.RegisterRepo(repoName, d); err != nil {
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
		info, err := gh.GetInfo(v)
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
			if err = t.ExecuteTemplate(buf, "upstreamUrl", &info); err != nil {
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
		if err = r.SetConfig(cfg); err != nil {
			return err
		}
	}
	return nil
}
