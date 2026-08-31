package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	satlogger "github.com/container-registry/harbor-satellite/internal/logger"
	proxyhandler "github.com/container-registry/harbor-satellite/internal/satellite/proxy"
	proxyimage "github.com/container-registry/harbor-satellite/internal/satellite/proxy/process/image"
	"github.com/container-registry/harbor-satellite/internal/satellite/store"
	"github.com/container-registry/harbor-satellite/pkg/config"
)

const (
	defaultListenAddress  = ":8585"
	defaultDataDir        = "/data/oci"
	defaultShutdownPeriod = 30 * time.Second
)

type options struct {
	listenAddress string
	dataDir       string
	upstream      string
	username      string
	password      string
	mode          proxyhandler.Mode
	plainHTTP     bool
	tlsCAFile     string
	tlsSkipVerify bool
	shutdown      time.Duration
}

func main() {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, _ = satlogger.InitLogger(ctx, "info", false, nil)

	if err := run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}

func parseOptions(arguments []string) (options, error) {
	mode, err := proxyhandler.ParseMode(envOrDefault("PROXY_MODE", proxyhandler.ModeProxy.String()))
	if err != nil {
		return options{}, err
	}
	shutdown, err := time.ParseDuration(envOrDefault("PROXY_SHUTDOWN_TIMEOUT", defaultShutdownPeriod.String()))
	if err != nil {
		return options{}, fmt.Errorf("parse PROXY_SHUTDOWN_TIMEOUT: %w", err)
	}

	opts := options{
		listenAddress: envOrDefault("PROXY_LISTEN_ADDRESS", defaultListenAddress),
		dataDir:       envOrDefault("PROXY_DATA_DIR", defaultDataDir),
		upstream:      os.Getenv("HARBOR_URL"),
		username:      os.Getenv("HARBOR_USERNAME"),
		password:      os.Getenv("HARBOR_PASSWORD"),
		mode:          mode,
		plainHTTP:     envBool("HARBOR_PLAIN_HTTP"),
		tlsCAFile:     os.Getenv("HARBOR_CA_FILE"),
		tlsSkipVerify: envBool("HARBOR_TLS_SKIP_VERIFY"),
		shutdown:      shutdown,
	}

	flags := flag.NewFlagSet("proxy", flag.ContinueOnError)
	flags.StringVar(&opts.listenAddress, "listen-address", opts.listenAddress, "HTTP listen address")
	flags.StringVar(&opts.dataDir, "data-dir", opts.dataDir, "OCI layout data directory")
	flags.StringVar(&opts.upstream, "upstream", opts.upstream, "upstream Harbor registry URL")
	flags.StringVar(&opts.username, "username", opts.username, "upstream registry username")
	flags.StringVar(&opts.password, "password", opts.password, "upstream registry password")
	flags.Var(&opts.mode, "mode", "operating mode: proxy or replica")
	flags.BoolVar(&opts.plainHTTP, "plain-http", opts.plainHTTP, "use HTTP for the upstream registry")
	flags.StringVar(&opts.tlsCAFile, "tls-ca-file", opts.tlsCAFile, "upstream registry CA certificate")
	flags.BoolVar(&opts.tlsSkipVerify, "tls-skip-verify", opts.tlsSkipVerify, "skip upstream TLS certificate verification")
	flags.DurationVar(&opts.shutdown, "shutdown-timeout", opts.shutdown, "graceful shutdown timeout")
	if err := flags.Parse(arguments); err != nil {
		return options{}, err
	}
	if err := opts.validate(); err != nil {
		return options{}, err
	}
	return opts, nil
}

func (o options) validate() error {
	if _, _, err := net.SplitHostPort(o.listenAddress); err != nil {
		return fmt.Errorf("invalid listen address %q: %w", o.listenAddress, err)
	}
	if o.dataDir == "" {
		return errors.New("OCI data directory is required")
	}
	if o.upstream == "" {
		return errors.New("upstream registry is required; set HARBOR_URL or --upstream")
	}
	if !o.mode.Valid() {
		return fmt.Errorf("invalid proxy mode %q", o.mode)
	}
	if o.shutdown <= 0 {
		return errors.New("shutdown timeout must be positive")
	}
	return nil
}

func run(ctx context.Context, opts options) error {
	localStore, err := store.NewOCIStore(opts.dataDir)
	if err != nil {
		return fmt.Errorf("initialize local OCI store: %w", err)
	}
	remoteStore, err := store.NewRegistryStore(store.RegistryOptions{
		Endpoint:  opts.upstream,
		Username:  opts.username,
		Password:  opts.password,
		PlainHTTP: opts.plainHTTP,
		TLS: config.TLSConfig{
			CAFile:     opts.tlsCAFile,
			SkipVerify: opts.tlsSkipVerify,
		},
	})
	if err != nil {
		return fmt.Errorf("initialize upstream registry store: %w", err)
	}

	listener, err := net.Listen("tcp", opts.listenAddress)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", opts.listenAddress, err)
	}

	server := &http.Server{
		Addr:              opts.listenAddress,
		Handler:           proxyhandler.New(proxyimage.NewPull(ctx, opts.mode, localStore, remoteStore)).Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       time.Minute,
	}
	serveResult := make(chan error, 1)
	go func() {
		err := server.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveResult <- err
	}()

	log.Printf("OCI proxy listening on %s in %s mode; upstream=%s data=%s", opts.listenAddress, opts.mode, opts.upstream, opts.dataDir)
	select {
	case err := <-serveResult:
		if err != nil {
			return fmt.Errorf("serve OCI proxy: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), opts.shutdown)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down OCI proxy: %w", err)
	}
	if err := <-serveResult; err != nil {
		return fmt.Errorf("serve OCI proxy: %w", err)
	}
	return nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envBool(name string) bool {
	value, err := strconv.ParseBool(os.Getenv(name))
	return err == nil && value
}
