package gitutil_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/gitutil"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/store/file"
	"os"
	"testing"
)

var gh *gitutil.GitHelper

const (
	remoteTemplate = `git@{{ .RepoHost }}:{{ .UpstreamOwner }}/{{ .Name }}`
)

func init() {
	l := zaplog.NewZapLogger(zaplog.DEBUG, "gitutil-test", true)
	l.InitLogger(nil)
	repoDb, err := file.NewDB(l, "git_repo_test.json", true)
	if err != nil {
		l.Fatalf("error creating repo database: %v", err)
	}
	gh = gitutil.NewHelper(repoDb, "alex.dhyatma@ubuntusoftware.net")
}

func TestGit(t *testing.T) {
	wd, _ := os.Getwd()
	repoInfo, err := gh.RepoInfo(wd)
	require.NoError(t, err)
	t.Log(repoInfo.String())
}
