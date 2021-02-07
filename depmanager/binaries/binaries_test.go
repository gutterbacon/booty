package binaries_test

import (
	"go.amplifyedge.org/booty-v2/depmanager/binaries"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaries(t *testing.T) {
	t.Run("testGrafana", testGrafana)
}

func testGrafana(t *testing.T) {
	gf := binaries.NewGrafana("7.4.0")
	err := gf.Download("./testdata")
	require.NoError(t, err)

	err = gf.Install()
	require.NoError(t, err)
}
