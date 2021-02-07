package fileutil_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/pkg/helpers/fileutil"
	"testing"
)

func TestFileUtil(t *testing.T) {
	t.Run("test copy file", testCopy)
}

func testCopy(t *testing.T) {
	err := fileutil.Copy("testdata/some_random_markdown_people_dont_care_about.md", "testdata/yeah.md")
	require.NoError(t, err)
}
