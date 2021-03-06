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

	err = downloader.Download("github.com/goreleaser/goreleaser?ref=v0.157.0", "./downloaded/goreleaser-0.157.0")
	require.NoError(t, err)

	err = downloader.Download("github.com/google/go-jsonnet?ref=v0.17.0", "./downloaded/jsonnet-0.17.0")
	require.NoError(t, err)
}
