package orchestrator_test

import (
	"testing"

	"go.amplifyedge.org/booty-v2/dep/orchestrator"

	"github.com/stretchr/testify/require"
)

var (
	composer *orchestrator.Orchestrator
)

func init() {
	composer = orchestrator.NewOrchestrator("booty")
}

func TestOrchestrator(t *testing.T) {
	t.Run("testDownloadAll", testDownloadAll)
	t.Run("testInstallSingle", testInstallSingle)
	t.Run("testInstallAll", testInstallAll)
	t.Run("testUninstallAll", testUninstallAll)
}

func testDownloadAll(t *testing.T) {
	err := composer.DownloadAll()
	require.NoError(t, err)
}

func testInstallAll(t *testing.T) {
	err := composer.InstallAll()
	require.NoError(t, err)
}

func testInstallSingle(t *testing.T) {
	err := composer.Install("goreleaser", "0.155.1")
	require.NoError(t, err)
}

func testUninstallAll(t *testing.T) {
	err := composer.UninstallAll()
	require.NoError(t, err)
}
