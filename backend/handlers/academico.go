package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"proyecto-grupo4/backend/models"
)

type DocenteStore interface {
	List(context.Context) ([]models.Docente, error)
	Create(context.Context, models.Docente) (models.Docente, error)
	Update(context.Context, models.Docente) error
	Delete(context.Context, int64) error
}

type DocentesHandler struct{ Store DocenteStore }

func (h DocentesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveCollection(w, r, "/api/docentes", func() { h.list(w, r) }, func() { h.create(w, r) }, func(id int64) { h.update(w, r, id) }, func(id int64) { h.delete(w, r, id) })
}

func (h DocentesHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeCatalogError(w, err, "docente no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h DocentesHandler) create(w http.ResponseWriter, r *http.Request) {
	var item models.Docente
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validateDocente(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), item)
	if err != nil {
		writeCatalogError(w, err, "docente no encontrado")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h DocentesHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var item models.Docente
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	item.ID = id
	if err := validateDocente(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Store.Update(r.Context(), item); err != nil {
		writeCatalogError(w, err, "docente no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h DocentesHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "docente no encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateDocente(item *models.Docente) error {
	item.Nombre = strings.TrimSpace(item.Nombre)
	item.Especialidad = strings.TrimSpace(item.Especialidad)
	item.Correo = strings.TrimSpace(item.Correo)
	if item.Nombre == "" || item.Especialidad == "" || item.Correo == "" {
		return validationError("nombre, especialidad y correo son obligatorios")
	}
	if _, err := mail.ParseAddress(item.Correo); err != nil {
		return validationError("correo inválido")
	}
	return nil
}

type ClaseStore interface {
	List(context.Context) ([]models.Clase, error)
	Create(context.Context, models.Clase) (models.Clase, error)
	Update(context.Context, models.Clase) error
	Delete(context.Context, int64) error
}

type ClasesHandler struct{ Store ClaseStore }

func (h ClasesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveCollection(w, r, "/api/clases", func() { h.list(w, r) }, func() { h.create(w, r) }, func(id int64) { h.update(w, r, id) }, func(id int64) { h.delete(w, r, id) })
}

func (h ClasesHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeCatalogError(w, err, "clase no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h ClasesHandler) create(w http.ResponseWriter, r *http.Request) {
	var item models.Clase
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validateClase(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), item)
	if err != nil {
		writeCatalogError(w, err, "docente no encontrado")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h ClasesHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var item models.Clase
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	item.ID = id
	if err := validateClase(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Store.Update(r.Context(), item); err != nil {
		writeCatalogError(w, err, "clase o docente no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h ClasesHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "clase no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateClase(item *models.Clase) error {
	item.Nombre = strings.TrimSpace(item.Nombre)
	item.Horario = strings.TrimSpace(item.Horario)
	if item.Nombre == "" || item.Horario == "" || item.DocenteID < 1 {
		return validationError("nombre, horario y docente son obligatorios")
	}
	return nil
}

type AsistenciaStore interface {
	List(context.Context) ([]models.Asistencia, error)
	Create(context.Context, models.Asistencia) (models.Asistencia, error)
	Update(context.Context, models.Asistencia) error
	Delete(context.Context, int64) error
}

type AsistenciasHandler struct{ Store AsistenciaStore }

func (h AsistenciasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveCollection(w, r, "/api/asistencias", func() { h.list(w, r) }, func() { h.create(w, r) }, func(id int64) { h.update(w, r, id) }, func(id int64) { h.delete(w, r, id) })
}

func (h AsistenciasHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeCatalogError(w, err, "asistencia no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h AsistenciasHandler) create(w http.ResponseWriter, r *http.Request) {
	var item models.Asistencia
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validateAsistencia(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), item)
	if err != nil {
		writeCatalogError(w, err, "matrícula o clase no encontrada")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h AsistenciasHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var item models.Asistencia
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	item.ID = id
	if err := validateAsistencia(&item); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Store.Update(r.Context(), item); err != nil {
		writeCatalogError(w, err, "asistencia, matrícula o clase no encontrada")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h AsistenciasHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "asistencia no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateAsistencia(item *models.Asistencia) error {
	item.Fecha = strings.TrimSpace(item.Fecha)
	if item.MatriculaID < 1 || item.ClaseID < 1 || !models.EstadosAsistencia[item.Estado] {
		return validationError("matrícula, clase y un estado válido son obligatorios")
	}
	if _, err := time.Parse("2006-01-02", item.Fecha); err != nil {
		return validationError("la fecha debe tener formato AAAA-MM-DD")
	}
	return nil
}

func serveCollection(w http.ResponseWriter, r *http.Request, collection string, list, create func(), update, delete func(int64)) {
	if r.URL.Path == collection {
		switch r.Method {
		case http.MethodGet:
			list()
		case http.MethodPost:
			create()
		default:
			methodNotAllowed(w)
		}
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "api" || parts[1] != strings.TrimPrefix(collection, "/api/") {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}
	switch r.Method {
	case http.MethodPut:
		update(id)
	case http.MethodDelete:
		delete(id)
	default:
		methodNotAllowed(w)
	}
}

type validationError string

func (e validationError) Error() string { return string(e) }
