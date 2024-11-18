package main

import (
	"ShoppingList/internal/config"
	"ShoppingList/internal/logger"
	"ShoppingList/internal/storage"
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var configPath string
	var secretsPath string
	flag.StringVar(&configPath, "config", "./../../config/local.yml", "specify path to config")
	flag.StringVar(&secretsPath, "secrets", "./../../.env", "specify path to secrets, only local env")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.MustLoad(configPath, secretsPath)

	logger.SetupLogger(cfg.Env, "./log/core.log")
	slog.Info("successfully loaded config and logger")
	slog.Debug("debug mode enabled")

	db, err := storage.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err.Error())
	}
	defer db.Close()

	done := make(chan struct{})
	go func() {
		// do some work
		close(done)
	}()

	select {
	case <-stop:
		slog.Info("external signal to stop tradesloader service", slog.Duration("shutdown timeout", 15*time.Second))
		cancel()
		select {
		case <-done:
			slog.Info("graceful shutdown done")
		case <-time.After(15 * time.Second):
			slog.Warn("shutdown timeout")
		}
	case <-done:
		slog.Warn("unexpected service shutdown")
	}
}
