package components_test

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/logging/zaplog"
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/store/file"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	db store.Storer
)

func init() {
	_ = os.Setenv("BOOTY_HOME", "./binaries-test")
	l := zaplog.NewZapLogger(zaplog.WARN, "store-test", true)
	l.InitLogger(nil)
	var err error
	db, err = file.NewDB(l, filepath.Join(osutil.GetDataDir(), "pkgs.json"))
	if err != nil {
		l.Fatal(err)
	}
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
	t.Run("testJb", testJb)
	t.Run("testVictoriaMetrics", testVictoriaMetrics)
	t.Run("testMkcert", testMkcert)
}

func testGrafana(t *testing.T) {
	var err error
	gf := components.NewGrafana(db)
	gf.SetVersion("7.4.2")
	err = gf.Download()
	require.NoError(t, err)

	err = gf.Install()
	require.NoError(t, err)

	err = gf.Uninstall()
	require.NoError(t, err)
}

func testGoreleaser(t *testing.T) {
	gor := components.NewGoreleaser(db)
	gor.SetVersion("0.155.1")
	err := gor.Download()
	require.NoError(t, err)

	// install
	err = gor.Install()
	require.NoError(t, err)

	// update
	err = gor.Update("0.156.2")
	require.NoError(t, err)

	// uninstall
	err = gor.Uninstall()
	require.NoError(t, err)
}

func testCaddy(t *testing.T) {
	cdy := components.NewCaddy(db)
	cdy.SetVersion("2.3.0")
	err := cdy.Download()
	require.NoError(t, err)

	// install
	err = cdy.Install()
	require.NoError(t, err)

	// uninstall
	err = cdy.Uninstall()
	require.NoError(t, err)
}

func testProto(t *testing.T) {
	pgc := components.NewProtocGenCobra(db)
	pgc.SetVersion("0.4.1")

	pgg := components.NewProtocGenGo(db)
	pgg.SetVersion("1.25.0")

	pgrpc := components.NewProtocGenGoGrpc(db)
	pgrpc.SetVersion("1.1.0")

	p := components.NewProtoc(db, []dep.Component{
		pgc, pgg, pgrpc,
	})
	p.SetVersion("3.13.0")
	err := p.Download()
	require.NoError(t, err)

	// install
	err = p.Install()
	require.NoError(t, err)

	// update
	err = p.Update("3.14.0")
	require.NoError(t, err)

	// uninstall
	err = p.Uninstall()
	require.NoError(t, err)

}

func testProtocGenGo(t *testing.T) {
	p := components.NewProtocGenGo(db)
	p.SetVersion("1.25.0")
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
	p := components.NewProtocGenCobra(db)
	p.SetVersion("0.4.1")
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
	p := components.NewProtocGenGoGrpc(db)
	p.SetVersion("1.1.0")
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
	g := components.NewGoJsonnet(db)
	g.SetVersion("0.17.0")
	err := g.Download()
	require.NoError(t, err)

	// install
	err = g.Install()
	require.NoError(t, err)

	//// update
	err = g.Update("0.17.0")
	require.NoError(t, err)

	//// uninstall
	err = g.Uninstall()
	require.NoError(t, err)
}

func testVictoriaMetrics(t *testing.T) {
	var err error
	g := components.NewVicMet(db)
	g.SetVersion("1.53.0")
	err = g.Download()
	require.NoError(t, err)

	// install
	err = g.Install()
	require.NoError(t, err)

	// update
	err = g.Update("1.53.1")
	require.NoError(t, err)

	// run
	err = g.Run()
	require.NoError(t, err)

	// // stop
	err = g.RunStop()
	require.NoError(t, err)

	// uninstall
	err = g.Uninstall()
	require.NoError(t, err)
}

func testJb(t *testing.T) {
	g := components.NewJb(db)
	g.SetVersion("0.4.0")
	err := g.Download()
	require.NoError(t, err)

	// install
	err = g.Install()
	require.NoError(t, err)

	// update
	err = g.Update("0.4.0")
	require.NoError(t, err)

	//// uninstall
	err = g.Uninstall()
	require.NoError(t, err)
}

func testMkcert(t *testing.T) {
	m := components.NewMkcert(db)
	m.SetVersion("1.4.3")
	err := m.Download()
	require.NoError(t, err)

	err = m.Install()
	require.NoError(t, err)

	err = m.Update("1.4.3")
	require.NoError(t, err)

	err = m.Uninstall()
	require.NoError(t, err)
}
