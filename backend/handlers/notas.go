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
	p := strings.TrimPrefix(r.URL.Path, "/api/notas/")

	// PATRÓN MEMENTO: Endpoint /undo para revertir la calificación al snapshot previo
	if strings.HasSuffix(p, "undo") {
		id, _ := strconv.ParseInt(strings.TrimSuffix(p, "/undo"), 10, 64)
		var mID int64
		var memento models.NotaMemento
		// Caretaker: busca en SQLite el último memento guardado de esta nota
		if err := h.DB.QueryRow("SELECT id, calificacion FROM mementos WHERE nota_id = ? ORDER BY id DESC LIMIT 1", id).Scan(&mID, &memento.Calificacion); err != nil {
			http.Error(w, "Sin historial previo para deshacer", 400); return
		}
		var n models.Nota
		h.DB.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&n.ID, &n.Estudiante, &n.Curso, &n.Calificacion)
		
		// Originator restaura su calificación interna desde el memento
		n.Restaurar(memento)
		h.DB.Exec("UPDATE notas SET calificacion = ? WHERE id = ?", n.Calificacion, id)
		h.DB.Exec("DELETE FROM mementos WHERE id = ?", mID) // Quita el memento usado
		json.NewEncoder(w).Encode(n); return
	}

	id, _ := strconv.ParseInt(p, 10, 64)
	switch r.Method {
	case "GET":
		rows, _ := h.DB.Query("SELECT id, estudiante, curso, calificacion FROM notas")
		defer rows.Close()
		res := []models.Nota{}
		for rows.Next() {
			var n models.Nota
			rows.Scan(&n.ID, &n.Estudiante, &n.Curso, &n.Calificacion)
			res = append(res, n)
		}
		json.NewEncoder(w).Encode(res)
	case "POST":
		var n models.Nota
		json.NewDecoder(r.Body).Decode(&n)
		
		// REGLA DE NEGOCIO: Solo se califica si el estudiante tiene matrícula 'Activa'
		var estado string
		err := h.DB.QueryRow("SELECT estado FROM matriculas WHERE estudiante = ? AND curso = ?", n.Estudiante, n.Curso).Scan(&estado)
		if err != nil || estado != "Activa" {
			http.Error(w, "No se puede calificar: la matrícula debe estar Activa", 400); return
		}
		res, _ := h.DB.Exec("INSERT INTO notas (estudiante, curso, calificacion) VALUES (?, ?, ?)", n.Estudiante, n.Curso, n.Calificacion)
		n.ID, _ = res.LastInsertId()
		json.NewEncoder(w).Encode(n)
	case "PUT":
		var n, act models.Nota
		json.NewDecoder(r.Body).Decode(&n)
		
		// PATRÓN MEMENTO: Antes de actualizar, Originator genera snapshot y Caretaker lo guarda en BD
		h.DB.QueryRow("SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?", id).Scan(&act.ID, &act.Estudiante, &act.Curso, &act.Calificacion)
		m := act.CrearMemento()
		h.DB.Exec("INSERT INTO mementos (nota_id, calificacion) VALUES (?, ?)", m.NotaID, m.Calificacion)

		// Actualiza con la nueva calificación
		h.DB.Exec("UPDATE notas SET estudiante = ?, curso = ?, calificacion = ? WHERE id = ?", n.Estudiante, n.Curso, n.Calificacion, id)
	case "DELETE":
		h.DB.Exec("DELETE FROM notas WHERE id = ?", id)
		h.DB.Exec("DELETE FROM mementos WHERE nota_id = ?", id)
		w.WriteHeader(204)
	}
}
