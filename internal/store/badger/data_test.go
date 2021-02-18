package badger_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/store/badger"
	"testing"
)

var (
	db  *badger.DB
	ips = []*store.InstalledPackage{
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
	l := zaplog.NewZapLogger(zaplog.DEBUG, "store-test", true)
	l.InitLogger(nil)
	db = badger.NewDB(l, "./testdata")
}

func TestStore(t *testing.T) {
	t.Run("testNew", testNew)
	t.Run("testGet", testGet)
	t.Run("testList", testList)
	t.Run("testDelete", testDelete)
}

func testNew(t *testing.T) {
	require.NoError(t, db.New(ips[0]))
	require.NoError(t, db.BulkNew(ips))
}

func testGet(t *testing.T) {
	queryName := "grafana"
	gf, err := db.Get(queryName)
	require.NoError(t, err)
	require.Equal(t, gf, ips[0])

	queryName = "bs-crypt"
	bc, err := db.Get(queryName)
	require.NoError(t, err)
	require.Equal(t, bc, ips[1])
}

func testList(t *testing.T) {
	// test when query is nil
	pkgs, err := db.List()
	require.NoError(t, err)
	require.Equal(t, len(ips), len(pkgs))
}

func testDelete(t *testing.T) {
	err := db.Delete("grafana")
	require.NoError(t, err)

	pkgs, _ := db.List()
	require.Equal(t, 1, len(pkgs))
}
