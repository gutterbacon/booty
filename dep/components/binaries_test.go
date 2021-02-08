package components_test

import (
	"os"

	"go.amplifyedge.org/booty-v2/dep/components"
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
	_ = os.MkdirAll("./testdata/db", 0755)
	db = store.NewDB(l, "./testdata/db")
}

func TestBinaries(t *testing.T) {
	t.Run("testGrafana", testGrafana)
	t.Run("testGoreleaser", testGoreleaser)
	t.Run("testCaddy", testCaddy)
	t.Run("testProto", testProto)
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

	// update
	err = cdy.Update("2.3.0")
	require.NoError(t, err)

	// uninstall
	err = cdy.Uninstall()
	require.NoError(t, err)
}

func testProto(t *testing.T) {
	p := components.NewProtoc(db, "3.13.0")
	err := p.Download("./testdata/downloads")
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	// update
	err = p.Update("3.14.0")
	require.NoError(t, err)

	// run
	err = p.Run("-I./proto/v2/", "-I.", "--go_out=./prototest/", "--go_opt=paths=source_relative", "./prototest/test.proto")
	require.NoError(t, err)
	_ = os.RemoveAll("./prototest/prototest")

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)

}
