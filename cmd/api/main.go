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
	// Comprobar si el archivo existe
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	
	_, err := fs.Stat(h.staticFS, path)
	if os.IsNotExist(err) {
		// Fallback a index.html
		r.URL.Path = "/"
	}

	h.fileServer.ServeHTTP(w, r)
}

func main() {
	fmt.Println("Facturacion Local - Iniciando Servidor")

	// Inicializar Base de Datos
	db, err := storage.InitDB("facturacion.db")
	if err != nil {
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}

	mux := http.NewServeMux()

	// 1. Configurar la API
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	// Registrar las rutas de la API (invoices, clients, etc)
	handlers.RegisterRoutes(mux, db)

	// 2. Configurar el Frontend (React SPA)
	// Extraer el subdirectorio "dist"
	distFS, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Error cargando archivos estáticos:", err)
	}

	spa := &spaHandler{
		staticFS:   distFS,
		fileServer: http.FileServer(http.FS(distFS)),
	}
	
	// Todo lo que no sea /api/ se delega a la SPA
	mux.Handle("/", spa)

	port := "8080"
	fmt.Printf("Servidor escuchando en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
