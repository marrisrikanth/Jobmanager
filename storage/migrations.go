package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	Version int
	UpSQL string
}

func applyMigrations(db *sql.DB, migrations []Migration) error {
	//Migration tracking table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)

	if err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	//Get applied migrations
	applied := make(map[int]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err == nil {
			applied[v] = true
		}
	}

	//Apply missing migrations in order
	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		fmt.Printf("Applying migration v%d...\n", m.Version)
		if _,err := db.Exec(m.UpSQL); err != nil {
			return fmt.Errorf("migration v%d failed: %w", m.Version, err)
		}
		_, err := db.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)", m.Version, time.Now())
		if err != nil {
			return fmt.Errorf("failed to record migration v%d: %w", m.Version, err)
		}
	}
	return nil
}


