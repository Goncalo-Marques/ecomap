package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goncalo-marques/ecomap/server/internal/config"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
	"github.com/goncalo-marques/ecomap/server/internal/service"
	"github.com/goncalo-marques/ecomap/server/internal/store"
	transporthttp "github.com/goncalo-marques/ecomap/server/internal/transport/http"
)

// Default configuration values.
const (
	defaultAddressHTTP = ":8080"
)

// Server metadata.
var (
	BuildGitHash   string
	BuildTimestamp string
	Hostname       string
)

func main() {
	var err error

	if len(Hostname) == 0 {
		Hostname, err = os.Hostname()
		if err != nil {
			logging.Logger.Error("main: failed to get hostname", logging.Error(err))
			return
		}
	}

	// Initialize logger with metadata attributes.
	slogHandler := slog.NewJSONHandler(os.Stdout, nil)
	logging.Init(slogHandler,
		slog.String(logging.BuildGitHash, BuildGitHash),
		slog.String(logging.BuildTimestamp, BuildTimestamp),
		slog.String(logging.Hostname, Hostname),
	)

	// Initialize base context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load service configuration.
	serviceConfig, err := config.Load()
	if err != nil {
		logging.Logger.ErrorContext(ctx, "main: failed to load service configuration", logging.Error(err))
		return
	}

	// Set up store.
	store, err := store.New(ctx, serviceConfig.Database)
	if err != nil {
		logging.Logger.ErrorContext(ctx, "main: failed to set up store", logging.Error(err))
		return
	}
	defer store.Close()

	// Set up service.
	service := service.New(store)

	// Handle signals.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Set up http server.
	addressHTTP := defaultAddressHTTP
	if len(serviceConfig.ServerHTTP.Address) != 0 {
		addressHTTP = serviceConfig.ServerHTTP.Address
	}

	handlerHTTP := transporthttp.New(service)

	serverHTTP := &http.Server{
		Addr:     addressHTTP,
		Handler:  handlerHTTP,
		ErrorLog: slog.NewLogLogger(logging.Logger.Handler(), slog.LevelError),
	}

	go func() {
		err := serverHTTP.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logging.Logger.InfoContext(ctx, "main: server http terminated")
		} else {
			logging.Logger.ErrorContext(ctx, "main: server http terminated unexpectedly", logging.Error(err))
		}

		sigs <- syscall.SIGQUIT
	}()

	logging.Logger.InfoContext(ctx, "main: server initialized")

	// Handle server shutdown.
	sig := <-sigs
	switch sig {
	case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
		if err := serverHTTP.Shutdown(ctx); err != nil {
			logging.Logger.ErrorContext(ctx, "main: failed to shutdown server http", logging.Error(err))
		}
	}

	logging.Logger.InfoContext(ctx, "main: server terminated")
}
