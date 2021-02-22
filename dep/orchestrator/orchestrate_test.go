package orchestrator_test

import (
	"fmt"
	"go.amplifyedge.org/booty-v2/internal/errutil"
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

	c = composer.Component("fragana")
	require.Equal(t, nil, c)

	t.Log("listing all components")
	comps := composer.AllComponents()
	require.NotEqual(t, nil, comps)

	t.Log("listing all commands")
	cmds := composer.Command()
	expectedName := "booty"
	require.Equal(t, cmds.Name(), expectedName)
	require.True(t, cmds.HasAvailableSubCommands())

	l := composer.Logger()
	require.NotEqual(t, nil, l)

	t.Log("downloading all components")
	err := composer.DownloadAll()
	require.NoError(t, err)

	t.Log("installing all components")
	err = composer.InstallAll()
	require.NoError(t, err)

	t.Log("installing single component")
	err = composer.Install("goreleaser", "0.157.0")
	require.NoError(t, err)

	t.Log("installing single non-existent component")
	err = composer.Install("nonexistent", "0.0.1")
	require.Error(t, errutil.New(errutil.ErrInvalidComponent, fmt.Errorf("name: %s, version: %s", "nonexistent", "0.0.1")), err)

	t.Log("running single component")
	err = composer.Run("goreleaser", "-v")
	require.NoError(t, err)

	t.Log("uninstall single component")
	err = composer.Uninstall("goreleaser")
	require.NoError(t, err)

	t.Log("uninstall single non-existent component")
	err = composer.Uninstall("nonexistent")
	require.Error(t, errutil.New(errutil.ErrUninstallComponent, fmt.Errorf("name: %s, err: no package of that name available", "nonexistent")), err)

	t.Log("listing all components")
	_, err = composer.AllInstalledComponents()
	require.NoError(t, err)

	t.Log("running check for updates")
	checker := composer.Checker()
	require.NotEqual(t, nil, checker)

	err = checker.CheckNewReleases()
	require.NoError(t, err)

	t.Log("running backup of single component config (if any)")
	err = composer.Backup("goreleaser")
	require.NoError(t, err)

	t.Log("running backup of all components config (if any)")
	err = composer.BackupAll()
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	t.Log("uninstall all components")
	err = composer.UninstallAll()
	require.NoError(t, err)
}
