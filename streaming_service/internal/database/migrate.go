package database

import (
	"database/sql"
	"fmt"
	"os"
)

func ApplySQLFile(db *sql.DB, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", path, err)
	}
	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("apply migration %s: %w", path, err)
	}
	return nil
}
