package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/enoquesousa/godeploy-platform/internal/platform/iox"
	"github.com/enoquesousa/godeploy-platform/internal/platform/sqlpool"
	"github.com/enoquesousa/godeploy-platform/internal/proxy"

	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	os.Exit(run(logger))
}

func run(logger *slog.Logger) int {
	dbPath := getenv("GODEPLOY_DB", "godeploy.db")
	addr := getenv("GODEPLOY_PROXY_ADDR", ":8080")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Error("open sqlite", slog.Any("err", err))
		return 1
	}
	sqlpool.ForSQLite(db)
	defer iox.Close(db)

	p, err := proxy.New(proxy.Config{
		Addr:         addr,
		DB:           db,
		PollInterval: time.Second,
		Logger:       logger,
	})
	if err != nil {
		logger.Error("configure proxy", slog.Any("err", err))
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("godeploy proxy listening", slog.String("addr", addr), slog.String("db", dbPath))
	if err := p.Run(ctx); err != nil {
		logger.Error("proxy stopped", slog.Any("err", err))
		return 1
	}
	return 0
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
