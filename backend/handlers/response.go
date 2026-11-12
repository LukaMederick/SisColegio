package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"proyecto-grupo4/backend/repository"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeCatalogError(w http.ResponseWriter, err error, missingMessage string) {
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, missingMessage)
		return
	}
	if errors.Is(err, repository.ErrConflict) {
		writeError(w, http.StatusConflict, "el registro está relacionado con otros datos o ya existe")
		return
	}
	writeError(w, http.StatusInternalServerError, "no se pudo completar la operación")
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "método no permitido")
}
