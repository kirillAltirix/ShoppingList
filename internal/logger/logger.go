package logger

import (
	"ShoppingList/internal/config"
	"log"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

func SetupLogger(env string, logFilePath string) {
	switch env {
	case config.EnvLocal:
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	case config.EnvDev:
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	case config.EnvProd:
		slog.SetDefault(slog.New(slog.NewJSONHandler(
			&lumberjack.Logger{
				Filename:   logFilePath,
				MaxSize:    500, //megabytes
				MaxBackups: 6,
				MaxAge:     30, //days
				Compress:   true,
			},
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)))
	default:
		log.Fatal("Incorrect env value")
	}
}
