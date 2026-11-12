package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"proyecto-grupo4/backend/handlers"
	"proyecto-grupo4/backend/models"
	"proyecto-grupo4/backend/repository"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "educacion.db"
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		log.Fatal(err)
	}
	if err := initializeDB(db); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	matriculas := handlers.MatriculasHandler{Store: repository.MatriculaRepository{DB: db}}
	notaRepo := repository.NotaRepository{DB: db}
	caretaker := models.NuevoNotaCaretaker(notaRepo)
	notas := handlers.NotasHandler{
		Store:         notaRepo,
		Caretaker:     caretaker,
		MatriculaStore: repository.MatriculaRepository{DB: db},
	}
	docentes := handlers.DocentesHandler{Store: repository.DocenteRepository{DB: db}}
	clases := handlers.ClasesHandler{Store: repository.ClaseRepository{DB: db}}
	asistencias := handlers.AsistenciasHandler{Store: repository.AsistenciaRepository{DB: db}}
	mux.Handle("/api/matriculas", matriculas)
	mux.Handle("/api/matriculas/", matriculas)
	mux.Handle("/api/notas", notas)
	mux.Handle("/api/notas/", notas)
	mux.Handle("/api/docentes", docentes)
	mux.Handle("/api/docentes/", docentes)
	mux.Handle("/api/clases", clases)
	mux.Handle("/api/clases/", clases)
	mux.Handle("/api/asistencias", asistencias)
	mux.Handle("/api/asistencias/", asistencias)

	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("API escuchando en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}

func initializeDB(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS matriculas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		estudiante TEXT NOT NULL,
		apellido TEXT NOT NULL DEFAULT '',
		grado_curso TEXT NOT NULL DEFAULT '',
		seccion TEXT NOT NULL DEFAULT '',
		fecha_inscripcion TEXT NOT NULL DEFAULT '',
		tutor_responsable TEXT NOT NULL DEFAULT '',
		estado TEXT NOT NULL CHECK (estado IN ('Pendiente', 'Activa', 'Anulada', 'Finalizada'))
	);
	CREATE TABLE IF NOT EXISTS notas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		estudiante TEXT NOT NULL,
		curso TEXT NOT NULL,
		calificacion REAL NOT NULL CHECK (calificacion >= 0 AND calificacion <= 20)
	);
	CREATE TABLE IF NOT EXISTS docentes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		especialidad TEXT NOT NULL,
		correo TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS clases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		horario TEXT NOT NULL,
		docente_id INTEGER NOT NULL REFERENCES docentes(id) ON DELETE RESTRICT
	);
	CREATE TABLE IF NOT EXISTS asistencias (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		matricula_id INTEGER NOT NULL REFERENCES matriculas(id) ON DELETE RESTRICT,
		clase_id INTEGER NOT NULL REFERENCES clases(id) ON DELETE RESTRICT,
		fecha TEXT NOT NULL,
		estado TEXT NOT NULL CHECK (estado IN ('Presente', 'Tardanza', 'Ausente')),
		UNIQUE (matricula_id, clase_id, fecha)
	);
	CREATE TABLE IF NOT EXISTS nota_versiones (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nota_id INTEGER NOT NULL REFERENCES notas(id) ON DELETE CASCADE,
		estudiante TEXT NOT NULL,
		curso TEXT NOT NULL,
		calificacion REAL NOT NULL,
		fecha_modificacion TEXT NOT NULL
	);`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	columns, err := db.Query(`PRAGMA table_info(matriculas)`) 
	if err != nil {
		return err
	}
	defer columns.Close()

	present := map[string]bool{}
	for columns.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		if err := columns.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		present[name] = true
	}
	if err := columns.Err(); err != nil {
		return err
	}

	for _, columnName := range []string{"apellido", "grado_curso", "seccion", "fecha_inscripcion", "tutor_responsable"} {
		if !present[columnName] {
			if _, err := db.Exec("ALTER TABLE matriculas ADD COLUMN " + columnName + " TEXT NOT NULL DEFAULT ''"); err != nil {
				return err
			}
		}
	}

	return nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" || origin == "" {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
