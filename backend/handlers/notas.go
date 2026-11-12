package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"proyecto-grupo4/backend/models"
)

type NotasHandler struct{ DB *sql.DB }

func (h NotasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/notas")

	if path == "" || path == "/" {
		if r.Method == http.MethodGet {
			rows, _ := h.DB.Query("SELECT id, estudiante, curso, calificacion FROM notas")
			defer rows.Close()
			res := []models.Nota{}
			for rows.Next() {
				var n models.Nota
				rows.Scan(&n.ID, &n.Estudiante, &n.Curso, &n.Calificacion)
				res = append(res, n)
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		if r.Method == http.MethodPost {
			var n models.Nota
			json.NewDecoder(r.Body).Decode(&n)
			res, _ := h.DB.Exec("INSERT INTO notas (estudiante, curso, calificacion) VALUES (?, ?, ?)", n.Estudiante, n.Curso, n.Calificacion)
			n.ID, _ = res.LastInsertId()
			json.NewEncoder(w).Encode(n)
			return
		}
	}

	trimmed := strings.TrimPrefix(path, "/")
	if strings.HasSuffix(trimmed, "/undo") {
		id, _ := strconv.ParseInt(strings.TrimSuffix(trimmed, "/undo"), 10, 64)
		// Patrón MEMENTO: Recuperar último snapshot desde Caretaker (BD)
		var mID int64
		var memento models.NotaMemento
		err := h.DB.QueryRow("SELECT id, calificacion FROM mementos WHERE nota_id = ? ORDER BY id DESC LIMIT 1", id).Scan(&mID, &memento.Calificacion)
		if err != nil {
			http.Error(w, "Sin historial previo", http.StatusBadRequest)
			return
		}
		var n models.Nota
		h.DB.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&n.ID, &n.Estudiante, &n.Curso, &n.Calificacion)
		n.Restaurar(memento)
		h.DB.Exec("UPDATE notas SET calificacion = ? WHERE id = ?", n.Calificacion, id)
		h.DB.Exec("DELETE FROM mementos WHERE id = ?", mID)
		json.NewEncoder(w).Encode(n)
		return
	}

	id, _ := strconv.ParseInt(trimmed, 10, 64)
	if r.Method == http.MethodDelete {
		h.DB.Exec("DELETE FROM notas WHERE id = ?", id)
		h.DB.Exec("DELETE FROM mementos WHERE nota_id = ?", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == http.MethodPut {
		var n models.Nota
		json.NewDecoder(r.Body).Decode(&n)

		// Patrón MEMENTO: Originator crea memento antes de modificar
		var actual models.Nota
		h.DB.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&actual.ID, &actual.Estudiante, &actual.Curso, &actual.Calificacion)
		memento := actual.CrearMemento()
		h.DB.Exec("INSERT INTO mementos (nota_id, calificacion) VALUES (?, ?)", memento.NotaID, memento.Calificacion)

		h.DB.Exec("UPDATE notas SET estudiante = ?, curso = ?, calificacion = ? WHERE id = ?", n.Estudiante, n.Curso, n.Calificacion, id)
		w.WriteHeader(http.StatusOK)
	}
}
