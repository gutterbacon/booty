package downloader_test

import (
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/downloader"
	"os"
	"testing"
)

func init() {
	_ = os.MkdirAll("downloaded", 0755)
}

func TestDownloader(t *testing.T) {
	t.Log("test download")
	err := downloader.Download("https://github.com/caddyserver/caddy/releases/download/v2.3.0/caddy_2.3.0_linux_arm64.tar.gz", "./downloaded/caddy-2.3.0")
	require.NoError(t, err)

	err = downloader.GitClone("https://github.com/goreleaser/goreleaser", "./downloaded/goreleaser-master", "")
	require.NoError(t, err)
}
