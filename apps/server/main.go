package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	port := env("PORT", "8080")
	databasePath := env("DATABASE_PATH", "data/wub.db")

	db, err := openDatabase(databasePath)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	if err := migrateDatabase(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.GET("/api/health", func(c echo.Context) error {
		if err := db.PingContext(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	server.POST("/api/answers", createAnswerHandler(db))
	server.GET("/api/answers", listAnswersHandler(db))

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
		server.Logger.Fatal(err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
