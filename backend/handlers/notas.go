package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"proyecto-grupo4/backend/models"
)

type NotaStore interface {
	List(ctx context.Context) ([]models.Nota, error)
	GetByID(ctx context.Context, id int64) (models.Nota, error)
	Create(ctx context.Context, item models.Nota) (models.Nota, error)
	Update(ctx context.Context, item models.Nota) error
	Delete(ctx context.Context, id int64) error
}

type NotasHandler struct {
	Store          NotaStore
	Caretaker      *models.NotaCaretaker
	MatriculaStore MatriculaStore
}

func (h NotasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/notas" {
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
	// /api/notas/{id}
	if len(parts) == 3 && parts[0] == "api" && parts[1] == "notas" {
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

	// Rutas de Memento: /api/notas/{id}/historial y /api/notas/{id}/restaurar
	if len(parts) == 4 && parts[0] == "api" && parts[1] == "notas" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || id < 1 {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		switch parts[3] {
		case "historial":
			if r.Method == http.MethodGet {
				h.historial(w, r, id)
			} else {
				methodNotAllowed(w)
			}
		case "restaurar":
			if r.Method == http.MethodPost {
				h.restaurar(w, r, id)
			} else {
				methodNotAllowed(w)
			}
		default:
			http.NotFound(w, r)
		}
		return
	}

	http.NotFound(w, r)
}

func (h NotasHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudieron listar las notas")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h NotasHandler) create(w http.ResponseWriter, r *http.Request) {
	var item models.Nota
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validateNota(item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validarEstudianteRegistrado(r.Context(), item.Estudiante); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), item)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo crear la nota")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h NotasHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var input models.Nota
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	input.ID = id
	if err := validateNota(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validarEstudianteRegistrado(r.Context(), input.Estudiante); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 1. Obtener el estado actual (Originator)
	originator, err := h.Store.GetByID(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	// 2. Memento: Guardar copia completa de la información actual antes de editarla
	if h.Caretaker != nil {
		memento := originator.CreateMemento()
		if _, err := h.Caretaker.GuardarSnapshot(memento); err != nil {
			writeError(w, http.StatusInternalServerError, "no se pudo registrar la versión histórica")
			return
		}
	}

	// 3. Aplicar las modificaciones solicitadas
	originator.Estudiante = strings.TrimSpace(input.Estudiante)
	originator.Curso = strings.TrimSpace(input.Curso)
	originator.Calificacion = input.Calificacion

	if err := h.Store.Update(r.Context(), originator); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, originator)
}

func (h NotasHandler) historial(w http.ResponseWriter, r *http.Request, id int64) {
	if h.Caretaker == nil {
		writeError(w, http.StatusInternalServerError, "administrador de versiones no disponible")
		return
	}
	// Comprobar que la nota exista
	if _, err := h.Store.GetByID(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	historial, err := h.Caretaker.ObtenerHistorial(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo recuperar el historial de versiones")
		return
	}
	writeJSON(w, http.StatusOK, historial)
}

func (h NotasHandler) restaurar(w http.ResponseWriter, r *http.Request, id int64) {
	var payload struct {
		VersionID int64 `json:"version_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.VersionID < 1 {
		writeError(w, http.StatusBadRequest, "identificador de versión inválido")
		return
	}
	if h.Caretaker == nil {
		writeError(w, http.StatusInternalServerError, "administrador de versiones no disponible")
		return
	}

	// 1. Recuperar el estado actual de la nota (Originator)
	originator, err := h.Store.GetByID(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	// 2. Recuperar la versión solicitada a través del Caretaker
	memento, err := h.Caretaker.ObtenerVersion(payload.VersionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "versión histórica no encontrada")
		return
	}
	if memento.NotaID != id {
		writeError(w, http.StatusBadRequest, "la versión no pertenece a esta calificación")
		return
	}

	// 3. Requisito C: "Antes de restaurar, guardar también el estado actual como una nueva versión para evitar perder información"
	mementoActual := originator.CreateMemento()
	if _, err := h.Caretaker.GuardarSnapshot(mementoActual); err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo respaldar la versión previa a la restauración")
		return
	}

	// 4. Originator restaura sus datos a partir del Memento
	originator.Restore(memento)

	// 5. Persistir en la base de datos sin duplicar registros de notas
	if err := h.Store.Update(r.Context(), originator); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, originator)
}

func (h NotasHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h NotasHandler) validarEstudianteRegistrado(ctx context.Context, estudiante string) error {
	if h.MatriculaStore == nil {
		return nil
	}

	items, err := h.MatriculaStore.List(ctx)
	if err != nil {
		return errors.New("no se pudo validar la matrícula del estudiante")
	}

	nombre := strings.TrimSpace(strings.ToLower(estudiante))
	for _, matricula := range items {
		if strings.TrimSpace(strings.ToLower(matricula.Estudiante)) == nombre && (matricula.Estado == "Activa" || matricula.Estado == "Finalizada") {
			return nil
		}
	}

	return errors.New("solo se pueden registrar notas para estudiantes con matrícula activa o finalizada")
}

func validateNota(item models.Nota) error {
	if strings.TrimSpace(item.Estudiante) == "" {
		return errors.New("estudiante es obligatorio")
	}
	if strings.TrimSpace(item.Curso) == "" {
		return errors.New("curso es obligatorio")
	}
	if item.Calificacion < 0 || item.Calificacion > 20 {
		return errors.New("calificacion debe estar entre 0 y 20")
	}
	return nil
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "nota no encontrada")
		return
	}
	writeError(w, http.StatusInternalServerError, "no se pudo procesar la nota")
}
