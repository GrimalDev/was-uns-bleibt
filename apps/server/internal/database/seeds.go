package database

import (
	"database/sql"
	"fmt"
)

func Seed(db *sql.DB) error {
	const seedKey = "fake_answers_v1"

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin database seed: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT OR IGNORE INTO database_seeds (seed_key) VALUES (?)`, seedKey)
	if err != nil {
		return fmt.Errorf("register database seed: %w", err)
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check database seed registration: %w", err)
	}
	if inserted == 0 {
		return tx.Commit()
	}

	seedAnswers := []struct {
		brainPartID int
		phrase      string
	}{
		{1, "Le jardin de Freiburg"},
		{1, "Un decathlon"},
		{2, "Ma maman avant son deuil"},
		{2, "A mon beau pere"},
		{3, "All the leave are brown"},
		{3, "Du bist ein guter man geworden"},
		{4, "Son sac a dos de velo"},
		{4, "Un mouvement de Catan"},
		{5, "La cafetiere a l'italienne"},
		{5, "Des pfannekuchen"},
	}
	for _, item := range seedAnswers {
		if _, err := tx.Exec(
			`INSERT INTO answers (brain_part_id, phrase) VALUES (?, ?)`,
			item.brainPartID,
			item.phrase,
		); err != nil {
			return fmt.Errorf("insert seed answer: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit database seed: %w", err)
	}

	return nil
}
