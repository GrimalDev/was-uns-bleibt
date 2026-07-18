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
		{1, "Yvette espiegle"},
		{1, "Ma maman avant son deuil"},
		{1, "De la guitarre"},
		{1, "Il lit sur le canape"},
		{2, "Son pull"},
		{2, "Ma guitarre"},
		{2, "Un callin"},
		{2, "Un pat pat sur ma tete"},
		{3, "All the leave are brown"},
		{3, "500 milles"},
		{3, "Du bist ein guter man geworden"},
		{3, "Lilibutz"},
		{4, "Des pfannekuchen"},
		{4, "Thomas nachtisch"},
		{4, "Une glace au camping"},
		{4, "Une foret-noire"},
		{5, "La cafetiere a l'italienne"},
		{5, "L'odeur de son atelier"},
		{5, "L'eau de mer"},
		{5, "L'herbe et le chlore"},
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
