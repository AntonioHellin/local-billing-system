// Package main serves the Local Billing System HTTP API and embedded React frontend.
package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"facturacion-local/frontend"
	"facturacion-local/internal/handlers"
	"facturacion-local/internal/storage"
)

type spaHandler struct {
	staticFS   fs.FS
	fileServer http.Handler
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if requested file exists in static FS
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	_, err := fs.Stat(h.staticFS, path)
	if os.IsNotExist(err) {
		// Fallback to index.html for client-side routing
		r.URL.Path = "/"
	}

	h.fileServer.ServeHTTP(w, r)
}

func main() {
	fmt.Println("Local Billing System - Starting Server")

	// Resolve database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "facturacion.db"
	}

	// Initialize SQLite Database
	db, err := storage.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	mux := http.NewServeMux()

	// 1. API Health endpoint
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Register API endpoints (clients, invoices, etc.)
	handlers.RegisterRoutes(mux, db)

	// 2. Embedded React SPA Frontend
	distFS, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Error loading embedded static assets:", err)
	}

	spa := &spaHandler{
		staticFS:   distFS,
		fileServer: http.FileServer(http.FS(distFS)),
	}

	// Delegate non-API routes to SPA handler
	mux.Handle("/", spa)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
