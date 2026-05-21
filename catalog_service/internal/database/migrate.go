package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ApplySQLFile(ctx context.Context, pool *pgxpool.Pool, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", path, err)
	}
	if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("apply migration %s: %w", path, err)
	}
	return nil
}
