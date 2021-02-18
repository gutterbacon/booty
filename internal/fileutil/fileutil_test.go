package fileutil_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/fileutil"
	"testing"
)

func TestFileUtil(t *testing.T) {
	t.Run("test copy file", testCopy)
}

func testCopy(t *testing.T) {
	sum, err := fileutil.Copy("testdata/some_random_markdown_people_dont_care_about.md", "testdata/yeah.md")
	require.NoError(t, err)
	t.Logf("SHA SUM: %s\n", sum)
	require.NotEqual(t, "", sum)
}
