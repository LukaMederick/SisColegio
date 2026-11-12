package repository

import (
	"context"
	"database/sql"

	"proyecto-grupo4/backend/models"
)

type MatriculaRepository struct{ DB *sql.DB }

func (r MatriculaRepository) List(ctx context.Context) ([]models.Matricula, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, estudiante, apellido, grado_curso, seccion, fecha_inscripcion, tutor_responsable, estado FROM matriculas ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Matricula, 0)
	for rows.Next() {
		var item models.Matricula
		if err := rows.Scan(&item.ID, &item.Estudiante, &item.Apellido, &item.GradoCurso, &item.Seccion, &item.FechaInscripcion, &item.TutorResponsable, &item.Estado); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r MatriculaRepository) Create(ctx context.Context, item models.Matricula) (models.Matricula, error) {
	result, err := r.DB.ExecContext(ctx, `INSERT INTO matriculas (estudiante, apellido, grado_curso, seccion, fecha_inscripcion, tutor_responsable, estado) VALUES (?, ?, ?, ?, ?, ?, ?)`, item.Estudiante, item.Apellido, item.GradoCurso, item.Seccion, item.FechaInscripcion, item.TutorResponsable, item.Estado)
	if err != nil {
		return models.Matricula{}, err
	}
	item.ID, err = result.LastInsertId()
	return item, err
}

func (r MatriculaRepository) Update(ctx context.Context, item models.Matricula) (models.Matricula, error) {
	result, err := r.DB.ExecContext(ctx, `UPDATE matriculas SET estudiante = ?, apellido = ?, grado_curso = ?, seccion = ?, fecha_inscripcion = ?, tutor_responsable = ? WHERE id = ?`, item.Estudiante, item.Apellido, item.GradoCurso, item.Seccion, item.FechaInscripcion, item.TutorResponsable, item.ID)
	if err != nil {
		return models.Matricula{}, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return models.Matricula{}, err
	}
	if n == 0 {
		return models.Matricula{}, sql.ErrNoRows
	}
	err = r.DB.QueryRowContext(ctx, `SELECT id, estudiante, apellido, grado_curso, seccion, fecha_inscripcion, tutor_responsable, estado FROM matriculas WHERE id = ?`, item.ID).Scan(&item.ID, &item.Estudiante, &item.Apellido, &item.GradoCurso, &item.Seccion, &item.FechaInscripcion, &item.TutorResponsable, &item.Estado)
	return item, err
}

func (r MatriculaRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM matriculas WHERE id = ?`, id)
	if err != nil {
		return translateConstraint(err)
	}
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r MatriculaRepository) UpdateEstado(ctx context.Context, id int64, estado string) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE matriculas SET estado = ? WHERE id = ?`, estado, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}
