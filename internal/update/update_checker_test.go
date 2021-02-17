package update_test

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/internal/store/badger"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/update"
)

var (
	checker *update.Checker
)

func init() {
	l := zaplog.NewZapLogger(zaplog.WARN, "update-test", true)
	l.InitLogger(nil)
	_ = os.MkdirAll("./testdata/db", 0755)
	db := badger.NewDB(l, "./testdata/db")
	checker = update.NewChecker(l, db, map[update.RepositoryURL]update.Version{
		update.RepositoryURL("https://github.com/grafana/grafana"):   update.Version("7.4.0"),
		update.RepositoryURL("https://github.com/caddyserver/caddy"): update.Version("2.3.0"),
	}, func(r update.RepositoryURL, v update.Version) error {
		fmt.Printf("YAY UPDATING %s to %s", r, v)
		return nil
	})
}

func TestUpdater(t *testing.T) {
	t.Run("testCheckRelease", testUpdateCheck)
	t.Run("testGetLatest", testGetLatest)
	t.Run("testFallbackScrape", testFallbackScrape)
}

func testUpdateCheck(t *testing.T) {
	err := checker.CheckNewReleases()
	require.NoError(t, err)
}

func testGetLatest(t *testing.T) {
	ver, err := update.GetLatestVersion("https://github.com/grpc/grpc-go")
	require.NoError(t, err)
	require.Equal(t, "1.1.0", ver)
}

func testFallbackScrape(t *testing.T) {
	ver, err := update.FallbackScrape("https://github.com/grpc/grpc-go")
	require.NoError(t, err)
	require.Equal(t, "cmd/protoc-gen-go-grpc/v1.1.0", ver)

	ver, err = update.FallbackScrape("https://github.com/amplify-edge/protoc-gen-cobra")
	require.NoError(t, err)
	require.Equal(t, "v0.4.1", ver)
}