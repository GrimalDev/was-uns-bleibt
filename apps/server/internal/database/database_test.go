package database

import "testing"

func TestMigrateDatabaseCreatesAnswersTable(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("run migrations twice: %v", err)
	}

	var migrationCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&migrationCount); err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if migrationCount != 2 {
		t.Fatalf("migration count = %d, want 2", migrationCount)
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
	if _, err := Open(""); err == nil {
		t.Fatal("Open(\"\") returned nil error")
	}
}

func TestSeedDatabaseRunsOnce(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	if err := Seed(db); err != nil {
		t.Fatalf("seed database: %v", err)
	}
	if err := Seed(db); err != nil {
		t.Fatalf("rerun database seed: %v", err)
	}

	var answerCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM answers`).Scan(&answerCount); err != nil {
		t.Fatalf("count seeded answers: %v", err)
	}
	if answerCount != 10 {
		t.Fatalf("answer count = %d, want 10", answerCount)
	}

	var seedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM database_seeds WHERE seed_key = 'fake_answers_v1'`).Scan(&seedCount); err != nil {
		t.Fatalf("count database seeds: %v", err)
	}
	if seedCount != 1 {
		t.Fatalf("seed count = %d, want 1", seedCount)
	}
}
