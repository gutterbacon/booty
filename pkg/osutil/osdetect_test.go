package osutil_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/pkg/osutil"
	"testing"
)

func TestOsDetectUtils(t *testing.T) {
	t.Run("testCurChown", testUserChown)
}

func testUserChown(t *testing.T) {
	err := osutil.CurUserChown("/tmp/shit.txt")
	require.NoError(t, err)
}
