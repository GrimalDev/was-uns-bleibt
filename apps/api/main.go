package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed web/*
var embeddedWeb embed.FS

func main() {
	port := env("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	staticFS, err := fs.Sub(embeddedWeb, "web")
	if err != nil {
		log.Fatalf("static fs setup failed: %v", err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/", spaHandler(fileServer, staticFS))

	addr := ":" + port
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func spaHandler(next http.Handler, staticFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		if _, err := fs.Stat(staticFS, path[1:]); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		index, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			http.Error(w, "index not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(index)
	})
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
