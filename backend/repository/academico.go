package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"proyecto-grupo4/backend/models"
)

var ErrConflict = errors.New("registro relacionado o duplicado")

func translateConstraint(err error) error {
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "constraint failed") {
		return ErrConflict
	}
	return err
}

type DocenteRepository struct{ DB *sql.DB }

func (r DocenteRepository) List(ctx context.Context) ([]models.Docente, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, nombre, especialidad, correo FROM docentes ORDER BY nombre, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Docente, 0)
	for rows.Next() {
		var item models.Docente
		if err := rows.Scan(&item.ID, &item.Nombre, &item.Especialidad, &item.Correo); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r DocenteRepository) Create(ctx context.Context, item models.Docente) (models.Docente, error) {
	result, err := r.DB.ExecContext(ctx, `INSERT INTO docentes (nombre, especialidad, correo) VALUES (?, ?, ?)`, item.Nombre, item.Especialidad, item.Correo)
	if err != nil {
		return models.Docente{}, translateConstraint(err)
	}
	item.ID, err = result.LastInsertId()
	return item, err
}

func (r DocenteRepository) Update(ctx context.Context, item models.Docente) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE docentes SET nombre = ?, especialidad = ?, correo = ? WHERE id = ?`, item.Nombre, item.Especialidad, item.Correo, item.ID)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

func (r DocenteRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM docentes WHERE id = ?`, id)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

type ClaseRepository struct{ DB *sql.DB }

func (r ClaseRepository) List(ctx context.Context) ([]models.Clase, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT c.id, c.nombre, c.horario, c.docente_id, d.nombre FROM clases c JOIN docentes d ON d.id = c.docente_id ORDER BY c.nombre, c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Clase, 0)
	for rows.Next() {
		var item models.Clase
		if err := rows.Scan(&item.ID, &item.Nombre, &item.Horario, &item.DocenteID, &item.DocenteNombre); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r ClaseRepository) Create(ctx context.Context, item models.Clase) (models.Clase, error) {
	result, err := r.DB.ExecContext(ctx, `INSERT INTO clases (nombre, horario, docente_id) VALUES (?, ?, ?)`, item.Nombre, item.Horario, item.DocenteID)
	if err != nil {
		return models.Clase{}, translateConstraint(err)
	}
	item.ID, err = result.LastInsertId()
	return item, err
}

func (r ClaseRepository) Update(ctx context.Context, item models.Clase) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE clases SET nombre = ?, horario = ?, docente_id = ? WHERE id = ?`, item.Nombre, item.Horario, item.DocenteID, item.ID)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

func (r ClaseRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM clases WHERE id = ?`, id)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

type AsistenciaRepository struct{ DB *sql.DB }

func (r AsistenciaRepository) List(ctx context.Context) ([]models.Asistencia, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT a.id, a.matricula_id, m.estudiante, a.clase_id, c.nombre, a.fecha, a.estado
		FROM asistencias a JOIN matriculas m ON m.id = a.matricula_id JOIN clases c ON c.id = a.clase_id
		ORDER BY a.fecha DESC, a.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Asistencia, 0)
	for rows.Next() {
		var item models.Asistencia
		if err := rows.Scan(&item.ID, &item.MatriculaID, &item.Estudiante, &item.ClaseID, &item.Clase, &item.Fecha, &item.Estado); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r AsistenciaRepository) Create(ctx context.Context, item models.Asistencia) (models.Asistencia, error) {
	result, err := r.DB.ExecContext(ctx, `INSERT INTO asistencias (matricula_id, clase_id, fecha, estado) VALUES (?, ?, ?, ?)`, item.MatriculaID, item.ClaseID, item.Fecha, item.Estado)
	if err != nil {
		return models.Asistencia{}, translateConstraint(err)
	}
	item.ID, err = result.LastInsertId()
	return item, err
}

func (r AsistenciaRepository) Update(ctx context.Context, item models.Asistencia) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE asistencias SET matricula_id = ?, clase_id = ?, fecha = ?, estado = ? WHERE id = ?`, item.MatriculaID, item.ClaseID, item.Fecha, item.Estado, item.ID)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

func (r AsistenciaRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM asistencias WHERE id = ?`, id)
	if err != nil {
		return translateConstraint(err)
	}
	return ensureUpdated(result)
}

func ensureUpdated(result sql.Result) error {
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}
