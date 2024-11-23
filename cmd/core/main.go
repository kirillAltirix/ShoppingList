package main

import (
	"ShoppingList/internal/api/httpv1"
	"ShoppingList/internal/config"
	"ShoppingList/internal/logger"
	"ShoppingList/internal/repository/postgre"
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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

	db, err := postgre.NewPostgreConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err.Error())
	}
	defer db.Close()

	r := mux.NewRouter()
	s := httpv1.Build(db)
	s.Register(r)
	srv := &http.Server{
		Handler: r,
		Addr:    cfg.Address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			slog.Error("error from http server", slog.String("error message", err.Error()))
		}
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
