package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin database migration: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`); err != nil {
		return fmt.Errorf("create schema migrations table: %w", err)
	}
	if _, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS database_seeds (
		seed_key TEXT PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`); err != nil {
		return fmt.Errorf("create database seeds table: %w", err)
	}

	var currentVersion int
	if err := tx.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion); err != nil {
		return fmt.Errorf("read schema migration version: %w", err)
	}

	for version, migration := range []string{
		`CREATE TABLE answers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			brain_part_id INTEGER NOT NULL CHECK (brain_part_id > 0),
			phrase TEXT NOT NULL CHECK (length(trim(phrase)) > 0),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX answers_brain_part_id_idx ON answers (brain_part_id);`,
		`CREATE UNIQUE INDEX answers_brain_part_phrase_idx ON answers (brain_part_id, phrase);`,
	} {
		migrationVersion := version + 1
		if migrationVersion <= currentVersion {
			continue
		}

		if _, err := tx.Exec(migration); err != nil {
			return fmt.Errorf("apply database migration %d: %w", migrationVersion, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, migrationVersion); err != nil {
			return fmt.Errorf("record database migration %d: %w", migrationVersion, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit database migration: %w", err)
	}

	return nil
}
