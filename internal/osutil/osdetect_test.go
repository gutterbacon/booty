package osutil_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"go.amplifyedge.org/booty-v2/internal/osutil"
)

func TestOsDetectUtils(t *testing.T) {
	t.Run("testCurChown", testUserChown)
	t.Run("testGetDirs", testGetDirs)
}

func testUserChown(t *testing.T) {
	err := osutil.CurUserChown("/tmp/shit.txt")
	require.NoError(t, err)
}

func testGetDirs(t *testing.T) {
	bin := osutil.GetBinDir()
	require.NotEqual(t, "", bin)
}
