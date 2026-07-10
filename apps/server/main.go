package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed web/*
var embeddedWeb embed.FS

func main() {
	port := env("PORT", "8080")

	staticFS, err := fs.Sub(embeddedWeb, "web")
	if err != nil {
		log.Fatalf("static fs setup failed: %v", err)
	}

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	assetHandler := http.FileServer(http.FS(staticFS))
	server.GET("/*", func(c echo.Context) error {
		requestPath := strings.TrimPrefix(c.Request().URL.Path, "/")
		if requestPath == "" {
			assetHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}

		if _, err := fs.Stat(staticFS, requestPath); err == nil {
			assetHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}

		return c.FileFS("index.html", http.FS(staticFS))
	})

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
