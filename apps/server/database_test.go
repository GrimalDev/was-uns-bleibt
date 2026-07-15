package main

import "testing"

func TestMigrateDatabaseCreatesAnswersTable(t *testing.T) {
	db, err := openDatabase(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := migrateDatabase(db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	if err := migrateDatabase(db); err != nil {
		t.Fatalf("run migrations twice: %v", err)
	}

	var migrationCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&migrationCount); err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("migration count = %d, want 1", migrationCount)
	}

	var tableName string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'answers'`).Scan(&tableName); err != nil {
		t.Fatalf("find answers table: %v", err)
	}
	if tableName != "answers" {
		t.Fatalf("table name = %q, want answers", tableName)
	}

	var createdAt string
	if err := db.QueryRow(`
		INSERT INTO answers (brain_part_id, phrase)
		VALUES (?, ?)
		RETURNING created_at`, 1, "A memory").Scan(&createdAt); err != nil {
		t.Fatalf("insert answer: %v", err)
	}
	if createdAt == "" {
		t.Fatal("created_at is empty")
	}
}

func TestOpenDatabaseRejectsEmptyPath(t *testing.T) {
	if _, err := openDatabase(""); err == nil {
		t.Fatal("openDatabase(\"\") returned nil error")
	}
}
