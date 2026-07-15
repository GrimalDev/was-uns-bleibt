package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func openDatabase(path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is empty")
	}

	if path != ":memory:" {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("resolve database path: %w", err)
		}
		path = absolutePath

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	dsn := sqliteDSN(path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func sqliteDSN(path string) string {
	if path == ":memory:" {
		return "file::memory:?mode=memory&cache=shared&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)"
	}

	query := url.Values{}
	query.Add("_pragma", "busy_timeout(5000)")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "foreign_keys(ON)")

	return (&url.URL{Scheme: "file", Path: path, RawQuery: query.Encode()}).String()
}

func migrateDatabase(db *sql.DB) error {
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

	var currentVersion int
	err = tx.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion)
	if err != nil {
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
