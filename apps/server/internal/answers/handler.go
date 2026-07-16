package answers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

const maxAnswerLength = 40

type createRequest struct {
	BrainPartID int    `json:"brain_part_id"`
	Phrase      string `json:"phrase"`
}

type answer struct {
	ID          int64  `json:"id"`
	BrainPartID int    `json:"brain_part_id"`
	Phrase      string `json:"phrase"`
	CreatedAt   string `json:"created_at"`
}

func CreateHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var request createRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		}

		request.Phrase = strings.TrimSpace(request.Phrase)
		if request.BrainPartID <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "brain_part_id must be positive"})
		}
		if request.Phrase == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "phrase must not be empty"})
		}
		if utf8.RuneCountInString(request.Phrase) > maxAnswerLength {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "phrase must be 40 characters or fewer"})
		}

		createdAnswer, err := insert(c, db, request)
		if err != nil {
			return fmt.Errorf("create answer: %w", err)
		}

		return c.JSON(http.StatusCreated, createdAnswer)
	}
}

func ListHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.QueryContext(c.Request().Context(), `
			SELECT id, brain_part_id, phrase, created_at
			FROM answers
			ORDER BY created_at ASC, id ASC`)
		if err != nil {
			return fmt.Errorf("list answers: %w", err)
		}
		defer rows.Close()

		items := make([]answer, 0)
		for rows.Next() {
			var item answer
			if err := rows.Scan(&item.ID, &item.BrainPartID, &item.Phrase, &item.CreatedAt); err != nil {
				return fmt.Errorf("scan answer: %w", err)
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("iterate answers: %w", err)
		}

		return c.JSON(http.StatusOK, items)
	}
}

func insert(c echo.Context, db *sql.DB, request createRequest) (answer, error) {
	var created answer
	err := db.QueryRowContext(c.Request().Context(), `
		INSERT INTO answers (brain_part_id, phrase)
		VALUES (?, ?)
		RETURNING id, brain_part_id, phrase, created_at`,
		request.BrainPartID,
		request.Phrase,
	).Scan(&created.ID, &created.BrainPartID, &created.Phrase, &created.CreatedAt)
	if err != nil {
		return answer{}, err
	}

	return created, nil
}
