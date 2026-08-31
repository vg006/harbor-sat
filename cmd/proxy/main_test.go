package main

import (
	"testing"
	"time"

	proxyhandler "github.com/container-registry/harbor-satellite/internal/satellite/proxy"
	"github.com/stretchr/testify/require"
)

func TestParseOptionsUsesEnvironmentDefaults(t *testing.T) {
	t.Setenv("HARBOR_URL", "https://harbor.example.test")
	t.Setenv("HARBOR_USERNAME", "robot$proxy")
	t.Setenv("HARBOR_PASSWORD", "secret")
	t.Setenv("HARBOR_TLS_SKIP_VERIFY", "true")
	t.Setenv("PROXY_MODE", "replica")
	t.Setenv("PROXY_LISTEN_ADDRESS", "127.0.0.1:18585")
	t.Setenv("PROXY_DATA_DIR", t.TempDir())
	t.Setenv("PROXY_SHUTDOWN_TIMEOUT", "5s")

	opts, err := parseOptions(nil)

	require.NoError(t, err)
	require.Equal(t, "https://harbor.example.test", opts.upstream)
	require.Equal(t, "robot$proxy", opts.username)
	require.Equal(t, "secret", opts.password)
	require.True(t, opts.tlsSkipVerify)
	require.Equal(t, proxyhandler.ModeReplica, opts.mode)
	require.Equal(t, "127.0.0.1:18585", opts.listenAddress)
	require.Equal(t, 5*time.Second, opts.shutdown)
}

func TestParseOptionsFlagsOverrideEnvironment(t *testing.T) {
	t.Setenv("HARBOR_URL", "https://ignored.example.test")
	t.Setenv("PROXY_MODE", "replica")

	opts, err := parseOptions([]string{
		"--upstream", "http://registry.example.test:5000",
		"--listen-address", ":5001",
		"--data-dir", t.TempDir(),
		"--mode", "proxy",
		"--plain-http",
	})

	require.NoError(t, err)
	require.Equal(t, "http://registry.example.test:5000", opts.upstream)
	require.Equal(t, ":5001", opts.listenAddress)
	require.Equal(t, proxyhandler.ModeProxy, opts.mode)
	require.True(t, opts.plainHTTP)
}

func TestParseOptionsRejectsInvalidConfiguration(t *testing.T) {
	t.Run("missing upstream", func(t *testing.T) {
		t.Setenv("HARBOR_URL", "")
		_, err := parseOptions(nil)
		require.ErrorContains(t, err, "upstream registry is required")
	})

	t.Run("invalid listen address", func(t *testing.T) {
		t.Setenv("HARBOR_URL", "https://harbor.example.test")
		_, err := parseOptions([]string{"--listen-address", "8585"})
		require.ErrorContains(t, err, "invalid listen address")
	})

	t.Run("invalid mode", func(t *testing.T) {
		t.Setenv("HARBOR_URL", "https://harbor.example.test")
		_, err := parseOptions([]string{"--mode", "cache"})
		require.ErrorContains(t, err, "invalid proxy mode")
	})
}
