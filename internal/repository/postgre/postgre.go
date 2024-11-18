package postgre

import (
	"ShoppingList/internal/config"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgreConnection(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open(
		"pq",
		fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to db with error: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db with error: %w", err)
	}

	return db, nil
}
