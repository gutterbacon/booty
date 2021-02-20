package orchestrator_test

import (
	"os"
	"testing"
	"time"

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
	c := composer.Component("grafana")
	require.NotEqual(t, nil, c)

	comps := composer.AllComponents()
	require.NotEqual(t, nil, comps)

	l := composer.Logger()
	require.NotEqual(t, nil, l)

	err := composer.DownloadAll()
	require.NoError(t, err)

	err = composer.InstallAll()
	require.NoError(t, err)

	err = composer.Install("goreleaser", "0.155.1")
	require.NoError(t, err)

	err = composer.Run("goreleaser", "-v")
	require.NoError(t, err)

	err = composer.Uninstall("goreleaser")
	require.NoError(t, err)

	_, err = composer.AllInstalledComponents()
	require.NoError(t, err)

	checker := composer.Checker()
	require.NotEqual(t, nil, checker)

	err = checker.CheckNewReleases()
	require.NoError(t, err)

	err = composer.Backup("goreleaser")
	require.NoError(t, err)

	err = composer.BackupAll()
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	err = composer.UninstallAll()
	require.NoError(t, err)
}
