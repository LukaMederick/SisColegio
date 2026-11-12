package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"proyecto-grupo4/backend/handlers"
)

func main() {
	db, err := sql.Open("sqlite", "educacion.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS matriculas (id INTEGER PRIMARY KEY AUTOINCREMENT, estudiante TEXT, curso TEXT, estado TEXT);
		CREATE TABLE IF NOT EXISTS notas (id INTEGER PRIMARY KEY AUTOINCREMENT, estudiante TEXT, curso TEXT, calificacion REAL);
		CREATE TABLE IF NOT EXISTS mementos (id INTEGER PRIMARY KEY AUTOINCREMENT, nota_id INTEGER, calificacion REAL);
	`)

	mux := http.NewServeMux()
	mux.Handle("/api/matriculas", handlers.MatriculasHandler{DB: db})
	mux.Handle("/api/matriculas/", handlers.MatriculasHandler{DB: db})
	mux.Handle("/api/notas", handlers.NotasHandler{DB: db})
	mux.Handle("/api/notas/", handlers.NotasHandler{DB: db})

	log.Println("API corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", cors(mux)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
