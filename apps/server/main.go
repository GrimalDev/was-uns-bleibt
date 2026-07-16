package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"was-uns-bleibt/server/internal/answers"
	"was-uns-bleibt/server/internal/database"
)

func main() {
	port := env("PORT", "8080")
	databasePath := env("DATABASE_PATH", "data/wub.db")

	db, err := database.Open(databasePath)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatalf("database seed failed: %v", err)
	}

	server := echo.New()
	hub := answers.NewHub()
	server.HideBanner = true
	server.HidePort = true

	server.GET("/api/health", func(c echo.Context) error {
		if err := db.PingContext(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	server.POST("/api/answers", answers.CreateHandler(db, hub))
	server.GET("/api/answers", answers.ListHandler(db))
	server.GET("/api/answers/ws", hub.Handler(db))

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
