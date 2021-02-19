package osutil_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"go.amplifyedge.org/booty-v2/internal/osutil"
)

func TestOsDetectUtils(t *testing.T) {
	t.Run("testGetDirs", testGetDirs)
}

func testGetDirs(t *testing.T) {
	bin := osutil.GetBinDir()
	require.NotEqual(t, "", bin)
}
