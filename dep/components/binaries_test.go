package components_test

import (
	"os"

	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/internal/osutil"

	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/store"

	"testing"

	"github.com/stretchr/testify/require"
)

var (
	db *store.DB
)

func init() {
	_ = os.Setenv("BOOTY_HOME", "./testdata")
	l := zaplog.NewZapLogger(zaplog.WARN, "store-test", true)
	l.InitLogger(nil)
	_ = os.MkdirAll("./testdata/db", 0755)
	db = store.NewDB(l, "./testdata/db")
}

func TestBinaries(t *testing.T) {
	t.Run("testGrafana", testGrafana)
	t.Run("testGoreleaser", testGoreleaser)
	t.Run("testCaddy", testCaddy)
	t.Run("testProtocGenGo", testProtocGenGo)
	t.Run("testProtocGenGrpc", testProtocGenGrpc)
	t.Run("testProtocGenCobra", testProtocGenCobra)
	t.Run("testProto", testProto)
	t.Run("testGoJsonnet", testGoJsonnet)
	t.Run("testVictoriaMetrics", testVictoriaMetrics)
}

func testGrafana(t *testing.T) {
	var err error
	gf := components.NewGrafana(db, "7.4.0")
	err = gf.Download()
	require.NoError(t, err)

	err = gf.Install()
	require.NoError(t, err)

	err = gf.Run()
	require.NoError(t, err)

	err = gf.RunStop()
	require.NoError(t, err)

	err = gf.Uninstall()
	require.NoError(t, err)
}

func testGoreleaser(t *testing.T) {
	gor := components.NewGoreleaser(db, "0.149.0")
	err := gor.Download()
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
	cdy := components.NewCaddy(db, "2.3.0")
	err := cdy.Download()
	require.NoError(t, err)

	// install
	err = cdy.Install()
	require.NoError(t, err)

	// update
	// err = cdy.Update("2.3.0")
	// require.NoError(t, err)

	// run
	err = cdy.Run()
	require.NoError(t, err)

	// stop
	err = cdy.RunStop()
	require.NoError(t, err)

	// uninstall
	err = cdy.Uninstall()
	require.NoError(t, err)
}

func testProto(t *testing.T) {
	p := components.NewProtoc(db, "3.13.0", []dep.Component{
		components.NewProtocGenCobra(db, "0.4.1"),
		components.NewProtocGenGo(db, "1.25.0"),
		components.NewProtocGenGoGrpc(db, "1.1.0"),
	})
	err := p.Download()
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	// update
	err = p.Update("3.14.0")
	require.NoError(t, err)

	// run
	err = p.Run("-I.", "--go_out=./prototest/", "--go_opt=paths=source_relative", "./prototest/test.proto")
	require.NoError(t, err)
	_ = os.RemoveAll("./prototest/prototest")

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)

}

func testProtocGenGo(t *testing.T) {
	p := components.NewProtocGenGo(db, "1.25.0")
	err := p.Download()
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	exists := osutil.ExeExists("protoc-gen-go")
	require.Equal(t, true, exists)

	// update
	err = p.Update("1.25.0")
	require.NoError(t, err)

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)
}

func testProtocGenCobra(t *testing.T) {
	p := components.NewProtocGenCobra(db, "0.4.1")
	err := p.Download()
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	exists := osutil.ExeExists("protoc-gen-cobra")
	require.Equal(t, true, exists)

	// update
	err = p.Update("0.4.1")
	require.NoError(t, err)

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)
}

func testProtocGenGrpc(t *testing.T) {
	p := components.NewProtocGenGoGrpc(db, "1.1.0")
	err := p.Download()
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	exists := osutil.ExeExists("protoc-gen-go-grpc")
	require.Equal(t, true, exists)

	// update
	err = p.Update("1.1.0")
	require.NoError(t, err)

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)
}

func testGoJsonnet(t *testing.T) {
	g := components.NewGoJsonnet(db, "0.17.0")
	err := g.Download()
	require.NoError(t, err)

	// install
	err = g.Install()
	require.NoError(t, err)

	// update
	err = g.Update("0.17.0")
	require.NoError(t, err)

	// uninstall
	err = g.Uninstall()
	require.NoError(t, err)
}

func testVictoriaMetrics(t *testing.T) {
	var err error
	g := components.NewVicMet(db, "1.53.0")
	//err = g.Download()
	//require.NoError(t, err)

	// install
	err = g.Install()
	require.NoError(t, err)

	// update
	// err = g.Update("1.53.0")
	// require.NoError(t, err)

	// run
	err = g.Run()
	require.NoError(t, err)

	// // stop
	err = g.RunStop()
	require.NoError(t, err)

	// uninstall
	//err = g.Uninstall()
	//require.NoError(t, err)
}
