package orchestrator_test

import (
	"os"
	"testing"

	"go.amplifyedge.org/booty-v2/dep/orchestrator"

	"github.com/stretchr/testify/require"
)

var (
	composer *orchestrator.Orchestrator
)

func init() {
	_ = os.MkdirAll("./testdata", 0755)
	_ = os.Setenv("BOOTY_HOME", "./testdata")
	composer = orchestrator.NewOrchestrator("booty")
}

func TestOrchestrator(t *testing.T) {
	t.Run("testAll", testAll)
}

func testAll(t *testing.T) {
	err := composer.DownloadAll()
	require.NoError(t, err)

	err = composer.InstallAll()
	require.NoError(t, err)

	err = composer.Install("goreleaser", "0.155.1")
	require.NoError(t, err)

	_, err = composer.AllInstalledComponents()
	require.NoError(t, err)

	err = composer.UninstallAll()
	require.NoError(t, err)
}
