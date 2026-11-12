package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Matricula struct {
	ID         int64  `json:"id"`
	Estudiante string `json:"estudiante"`
	Curso      string `json:"curso"`
	Estado     string `json:"estado"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "educacion.db")
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
	mux.HandleFunc("/api/matriculas", handleMatriculas)
	mux.HandleFunc("/api/matriculas/", handleMatriculaID)
	mux.HandleFunc("/api/notas", handleNotas)
	mux.HandleFunc("/api/notas/", handleNotaID)

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

func handleMatriculas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		rows, _ := db.Query("SELECT id, estudiante, curso, estado FROM matriculas")
		defer rows.Close()
		res := []Matricula{}
		for rows.Next() {
			var m Matricula
			rows.Scan(&m.ID, &m.Estudiante, &m.Curso, &m.Estado)
			res = append(res, m)
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	if r.Method == "POST" {
		var m Matricula
		json.NewDecoder(r.Body).Decode(&m)
		if m.Estado == "" {
			m.Estado = "Pendiente"
		}
		res, _ := db.Exec("INSERT INTO matriculas (estudiante, curso, estado) VALUES (?, ?, ?)", m.Estudiante, m.Curso, m.Estado)
		m.ID, _ = res.LastInsertId()
		json.NewEncoder(w).Encode(m)
	}
}

func handleMatriculaID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/matriculas/"), 10, 64)
	if r.Method == "DELETE" {
		db.Exec("DELETE FROM matriculas WHERE id = ?", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == "PUT" {
		var body struct {
			Estado string `json:"estado"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		var actual string
		db.QueryRow("SELECT estado FROM matriculas WHERE id = ?", id).Scan(&actual)

		// Patrón STATE: delegar validación al objeto de Estado
		if !ObtenerEstado(actual).PuedePasar(body.Estado) {
			http.Error(w, fmt.Sprintf("Transición inválida de %s a %s", actual, body.Estado), http.StatusBadRequest)
			return
		}
		db.Exec("UPDATE matriculas SET estado = ? WHERE id = ?", body.Estado, id)
		w.WriteHeader(http.StatusOK)
	}
}

func handleNotas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		rows, _ := db.Query("SELECT id, estudiante, curso, calificacion FROM notas")
		defer rows.Close()
		res := []Nota{}
		for rows.Next() {
			var n Nota
			rows.Scan(&n.ID, &n.Estudiante, &n.Curso, &n.Calificacion)
			res = append(res, n)
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	if r.Method == "POST" {
		var n Nota
		json.NewDecoder(r.Body).Decode(&n)
		res, _ := db.Exec("INSERT INTO notas (estudiante, curso, calificacion) VALUES (?, ?, ?)", n.Estudiante, n.Curso, n.Calificacion)
		n.ID, _ = res.LastInsertId()
		json.NewEncoder(w).Encode(n)
	}
}

func handleNotaID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notas/")
	if strings.HasSuffix(path, "/undo") {
		id, _ := strconv.ParseInt(strings.TrimSuffix(path, "/undo"), 10, 64)
		// Patrón MEMENTO: Caretaker recupera el último snapshot
		var mID int64
		var memento NotaMemento
		err := db.QueryRow("SELECT id, calificacion FROM mementos WHERE nota_id = ? ORDER BY id DESC LIMIT 1", id).Scan(&mID, &memento.Calificacion)
		if err != nil {
			http.Error(w, "Sin historial para deshacer", http.StatusBadRequest)
			return
		}
		var nota Nota
		db.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&nota.ID, &nota.Estudiante, &nota.Curso, &nota.Calificacion)
		nota.Restaurar(memento)
		db.Exec("UPDATE notas SET calificacion = ? WHERE id = ?", nota.Calificacion, id)
		db.Exec("DELETE FROM mementos WHERE id = ?", mID)
		json.NewEncoder(w).Encode(nota)
		return
	}

	id, _ := strconv.ParseInt(path, 10, 64)
	if r.Method == "DELETE" {
		db.Exec("DELETE FROM notas WHERE id = ?", id)
		db.Exec("DELETE FROM mementos WHERE nota_id = ?", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == "PUT" {
		var n Nota
		json.NewDecoder(r.Body).Decode(&n)
		// Patrón MEMENTO: Originator crea memento antes de mutar
		var actual Nota
		db.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&actual.ID, &actual.Estudiante, &actual.Curso, &actual.Calificacion)
		memento := actual.CrearMemento()
		db.Exec("INSERT INTO mementos (nota_id, calificacion) VALUES (?, ?)", memento.NotaID, memento.Calificacion)

		db.Exec("UPDATE notas SET estudiante = ?, curso = ?, calificacion = ? WHERE id = ?", n.Estudiante, n.Curso, n.Calificacion, id)
		w.WriteHeader(http.StatusOK)
	}
}
