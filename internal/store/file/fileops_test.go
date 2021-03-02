package file_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/store"
	"testing"

	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/store/file"
)

var (
	filedb *file.DB
	repoDb *file.DB
	ips    = []*store.InstalledPackage{
		{
			Name:    "grafana",
			Version: "7.4.0",
			FilesMap: map[string]string{
				"/usr/local/bin/grafana-server": "some_hash",
			},
		},
		{
			Name:    "bs-crypt",
			Version: "0.0.1",
			FilesMap: map[string]string{
				"/usr/local/bin/bs-crypt": "some_other_hash",
			},
		},
	}
	repos = map[string]string{
		"sys-share": "./testdata/sys-share",
		"sys":       "./testdata/sys",
	}
)

func init() {
	var err error
	l := zaplog.NewZapLogger(zaplog.DEBUG, "store-test", true)
	l.InitLogger(nil)
	filedb, err = file.NewDB(l, "./fileops_test.json", false)
	if err != nil {
		l.Fatalf("error creating file: %v", err)
	}
	repoDb, err = file.NewDB(l, "repoops_test.json", true)
	if err != nil {
		l.Fatalf("error creating repo database: %v", err)
	}
}

func TestFileDB(t *testing.T) {
	t.Run("testAll", testAll)
	t.Run("testRepo", testRepo)
}

func testAll(t *testing.T) {
	err := filedb.New(ips[0])
	require.NoError(t, err)

	err = filedb.New(ips[1])
	require.NoError(t, err)

	newIp := &store.InstalledPackage{
		Name:    "bs-crypt",
		Version: "0.0.2",
		FilesMap: map[string]string{
			"/usr/local/bin/bs-crypt": "some_hash",
		},
	}
	err = filedb.New(newIp)
	require.NoError(t, err)

	ip, err := filedb.Get("grafana")
	require.NoError(t, err)
	require.Equal(t, ips[0], ip)

	listIps, err := filedb.List()
	require.NoError(t, err)
	require.Equal(t, len(ips), len(listIps))

	err = filedb.Delete("grafana")
	require.NoError(t, err)

	listIps, _ = filedb.List()
	require.Equal(t, 1, len(listIps))
	require.Equal(t, newIp, listIps[0])
}

func testRepo(t *testing.T) {
	for k, v := range repos {
		err := repoDb.RegisterRepo(k, v)
		require.NoError(t, err)

		dirPath := repoDb.GetRepo(k)
		require.Equal(t, v, dirPath)
	}

	listRepo, err := repoDb.ListRepo()
	require.NoError(t, err)
	require.Equal(t, repos, listRepo)

	err = repoDb.UnregisterRepo("sys-share")
	require.NoError(t, err)

	listRepo, err = repoDb.ListRepo()
	require.Equal(t, 1, len(listRepo))
	require.Equal(t, repos["sys"], listRepo["sys"])

	err = repoDb.UnregisterAll()
	require.NoError(t, err)
}
