package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"proyecto-grupo4/backend/models"
)

type MatriculaStore interface {
	List(ctx context.Context) ([]models.Matricula, error)
	Create(ctx context.Context, item models.Matricula) (models.Matricula, error)
	Update(ctx context.Context, item models.Matricula) (models.Matricula, error)
	Delete(ctx context.Context, id int64) error
	UpdateEstado(ctx context.Context, id int64, estado string) error
}

type MatriculasHandler struct{ Store MatriculaStore }

func (h MatriculasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/matriculas" {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			methodNotAllowed(w)
		}
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 3 && parts[0] == "api" && parts[1] == "matriculas" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || id < 1 {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}
		switch r.Method {
		case http.MethodPut:
			h.update(w, r, id)
		case http.MethodDelete:
			h.delete(w, r, id)
		default:
			methodNotAllowed(w)
		}
		return
	}
	if len(parts) != 4 || parts[0] != "api" || parts[1] != "matriculas" || parts[3] != "estado" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPut {
		methodNotAllowed(w)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}
	var input struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo validar la matrícula")
		return
	}
	var matriculaEncontrada *models.Matricula
	for _, item := range items {
		if item.ID == id {
			matriculaEncontrada = &item
			break
		}
	}
	if matriculaEncontrada == nil {
		writeError(w, http.StatusNotFound, "matrícula no encontrada")
		return
	}
	if err := matriculaEncontrada.CambiarEstado(input.Estado); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err := h.Store.UpdateEstado(r.Context(), id, matriculaEncontrada.Estado); err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo actualizar la matrícula")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "estado actualizado", "estado": matriculaEncontrada.Estado})
}

func (h MatriculasHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var item models.Matricula
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	item.ID = id
	item.Estudiante = strings.TrimSpace(item.Estudiante)
	item.Apellido = strings.TrimSpace(item.Apellido)
	item.GradoCurso = strings.TrimSpace(item.GradoCurso)
	item.Seccion = strings.TrimSpace(item.Seccion)
	item.FechaInscripcion = strings.TrimSpace(item.FechaInscripcion)
	item.TutorResponsable = strings.TrimSpace(item.TutorResponsable)
	if item.Estudiante == "" || item.Apellido == "" || item.GradoCurso == "" || item.Seccion == "" || item.FechaInscripcion == "" || item.TutorResponsable == "" {
		writeError(w, http.StatusBadRequest, "nombre, apellido, grado/curso, sección, fecha de inscripción y tutor son obligatorios")
		return
	}
	updated, err := h.Store.Update(r.Context(), item)
	if err != nil {
		writeCatalogError(w, err, "matrícula no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h MatriculasHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "matrícula no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h MatriculasHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudieron listar las matrículas")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h MatriculasHandler) create(w http.ResponseWriter, r *http.Request) {
	var input models.Matricula
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	input.Estudiante = strings.TrimSpace(input.Estudiante)
	input.Apellido = strings.TrimSpace(input.Apellido)
	input.GradoCurso = strings.TrimSpace(input.GradoCurso)
	input.Seccion = strings.TrimSpace(input.Seccion)
	input.FechaInscripcion = strings.TrimSpace(input.FechaInscripcion)
	input.TutorResponsable = strings.TrimSpace(input.TutorResponsable)
	if input.Estudiante == "" || input.Apellido == "" || input.GradoCurso == "" || input.Seccion == "" || input.FechaInscripcion == "" || input.TutorResponsable == "" {
		writeError(w, http.StatusBadRequest, "nombre, apellido, grado/curso, sección, fecha de inscripción y tutor son obligatorios")
		return
	}
	if input.Estado == "" {
		input.Estado = "Activa"
	}
	if !models.EstadosMatricula[input.Estado] {
		writeError(w, http.StatusBadRequest, "estado inválido")
		return
	}
	created, err := h.Store.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo crear la matrícula")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
