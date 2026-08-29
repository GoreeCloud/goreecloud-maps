package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/auth"
	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/httpapi"
	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	listenAddress := envOrDefault("MAPS_LISTEN_ADDRESS", ":8080")
	databaseURL := os.Getenv("MAPS_DATABASE_URL")
	issuerURL := os.Getenv("MAPS_OIDC_ISSUER_URL")
	clientID := os.Getenv("MAPS_OIDC_CLIENT_ID")

	if databaseURL == "" || issuerURL == "" || clientID == "" {
		logger.Error("required configuration is missing",
			"requires", []string{"MAPS_DATABASE_URL", "MAPS_OIDC_ISSUER_URL", "MAPS_OIDC_CLIENT_ID"},
		)
		os.Exit(2)
	}

	dataStore, err := store.Open(ctx, databaseURL)
	if err != nil {
		logger.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	defer dataStore.Close()

	verifier, err := auth.NewVerifier(ctx, issuerURL, clientID)
	if err != nil {
		logger.Error("OIDC verifier initialization failed", "error", err)
		os.Exit(1)
	}

	handler := httpapi.New(logger, dataStore, verifier)
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownContext, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownContext.Done()
		graceContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(graceContext); err != nil {
			logger.Error("HTTP shutdown failed", "error", err)
		}
	}()

	logger.Info("GoreeCloud Maps API starting", "address", listenAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
