package update_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
)

var (
	checker *update.Checker
)

func init() {
	l := zaplog.NewZapLogger(zaplog.WARN, "update-test", true)
	l.InitLogger(nil)
	_ = os.MkdirAll("./testdata/db", 0755)
	db := store.NewDB(l, "./testdata/db")
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
}

func testUpdateCheck(t *testing.T) {
	err := checker.CheckNewReleases()
	require.NoError(t, err)
}

