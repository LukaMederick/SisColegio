package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"proyecto-grupo4/backend/models"
)

type MatriculasHandler struct{ DB *sql.DB }

func (h MatriculasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/matriculas")

	if path == "" || path == "/" {
		if r.Method == http.MethodGet {
			rows, _ := h.DB.Query("SELECT id, estudiante, curso, estado FROM matriculas")
			defer rows.Close()
			res := []models.Matricula{}
			for rows.Next() {
				var m models.Matricula
				rows.Scan(&m.ID, &m.Estudiante, &m.Curso, &m.Estado)
				res = append(res, m)
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		if r.Method == http.MethodPost {
			var m models.Matricula
			json.NewDecoder(r.Body).Decode(&m)
			if m.Estado == "" { m.Estado = "Pendiente" }
			res, _ := h.DB.Exec("INSERT INTO matriculas (estudiante, curso, estado) VALUES (?, ?, ?)", m.Estudiante, m.Curso, m.Estado)
			m.ID, _ = res.LastInsertId()
			json.NewEncoder(w).Encode(m)
			return
		}
	}

	id, _ := strconv.ParseInt(strings.TrimPrefix(path, "/"), 10, 64)
	if r.Method == http.MethodDelete {
		h.DB.Exec("DELETE FROM matriculas WHERE id = ?", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == http.MethodPut {
		var body struct{ Estado string `json:"estado"` }
		json.NewDecoder(r.Body).Decode(&body)

		var m models.Matricula
		h.DB.QueryRow("SELECT id, estudiante, curso, estado FROM matriculas WHERE id = ?", id).Scan(&m.ID, &m.Estudiante, &m.Curso, &m.Estado)

		// Aplicación del Patrón STATE
		if err := m.CambiarEstado(body.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.DB.Exec("UPDATE matriculas SET estado = ? WHERE id = ?", m.Estado, id)
		w.WriteHeader(http.StatusOK)
	}
}
