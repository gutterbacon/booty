package osdetect_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/pkg/helpers/osdetect"
	"testing"
)

func TestOsDetectUtils(t *testing.T) {
	t.Run("testCurChown", testUserChown)
}

func testUserChown(t *testing.T) {
	err := osdetect.CurUserChown("/tmp/shit.txt")
	require.NoError(t, err)
}
