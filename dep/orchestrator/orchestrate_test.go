package orchestrator_test

import (
	"go.amplifyedge.org/booty-v2/internal/osutil"
	"go.amplifyedge.org/booty-v2/internal/testutil"
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

	if osutil.GetOS() == "linux" || osutil.GetOS() == "darwin" {
		status := composer.Serve()
		//_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		_ = testutil.Kill()
		require.Equal(t, 0, status)
	}

	err = composer.Backup("goreleaser")
	require.NoError(t, err)

	err = composer.BackupAll()
	require.NoError(t, err)

	err = composer.UninstallAll()
	require.NoError(t, err)
}
