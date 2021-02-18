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
	ips    = []*store.InstalledPackage{
		{
			Name:    "grafana",
			Version: "7.4.0",
			FilesMap: map[string]int{
				"/usr/local/bin/grafana-server": 0755,
			},
		},
		{
			Name:    "bs-crypt",
			Version: "some_hash",
			FilesMap: map[string]int{
				"/usr/local/bin/bs-crypt": 0755,
			},
		},
	}
)

func init() {
	var err error
	l := zaplog.NewZapLogger(zaplog.DEBUG, "store-test", true)
	l.InitLogger(nil)
	filedb, err = file.NewDB(l, "./fileops_test.json")
	if err != nil {
		l.Fatalf("error creating file: %v", err)
	}
}

func TestFileDB(t *testing.T) {
	t.Run("testAll", testAll)
}

func testAll(t *testing.T) {
	err := filedb.New(ips[0])
	require.NoError(t, err)

	err = filedb.New(ips[1])
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
	require.Equal(t, ips[1], listIps[0])
}
