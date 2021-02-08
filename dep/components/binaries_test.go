package components_test

import (
	"go.amplifyedge.org/booty-v2/depmanager/components"
	"go.amplifyedge.org/booty-v2/pkg/logging/zaplog"
	"go.amplifyedge.org/booty-v2/pkg/store"

	"testing"

	"github.com/stretchr/testify/require"
)

var (
	db *store.DB
)

func init() {
	l := zaplog.NewZapLogger(zaplog.DEBUG, "store-test", true)
	l.InitLogger(nil)
	db = store.NewDB(l, "./testdata/db")
}

func TestBinaries(t *testing.T) {
	t.Run("testGrafana", testGrafana)
	t.Run("testGoreleaser", testGoreleaser)
	t.Run("testCaddy", testCaddy)
}

func testGrafana(t *testing.T) {
	gf := components.NewGrafana(db, "7.3.7")
	err := gf.Download("./testdata/downloads")
	require.NoError(t, err)

	err = gf.Install()
	require.NoError(t, err)

	err = gf.Update("7.4.0")
	require.NoError(t, err)

	err = gf.Uninstall()
	require.NoError(t, err)
}

func testGoreleaser(t *testing.T) {
	gor := components.NewGoreleaser(db, "0.149.0")
	err := gor.Download("./testdata/downloads")
	require.NoError(t, err)

	// install
	err = gor.Install()
	require.NoError(t, err)

	// update
	err = gor.Update("0.155.0")
	require.NoError(t, err)

	// uninstall
	err = gor.Uninstall()
	require.NoError(t, err)
}

func testCaddy(t *testing.T) {
	cdy := components.NewCaddy(db, "2.2.0")
	err := cdy.Download("./testdata/downloads")
	require.NoError(t, err)

	// install
	err = cdy.Install()
	require.NoError(t, err)

	//// update
	err = cdy.Update("2.3.0")
	require.NoError(t, err)

	//// uninstall
	err = cdy.Uninstall()
	require.NoError(t, err)
}
