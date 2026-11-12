package repository

import (
	"context"
	"database/sql"

	"proyecto-grupo4/backend/models"
)

type NotaRepository struct{ DB *sql.DB }

func (r NotaRepository) List(ctx context.Context) ([]models.Nota, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, estudiante, curso, calificacion FROM notas ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Nota, 0)
	for rows.Next() {
		var item models.Nota
		if err := rows.Scan(&item.ID, &item.Estudiante, &item.Curso, &item.Calificacion); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r NotaRepository) GetByID(ctx context.Context, id int64) (models.Nota, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT id, estudiante, curso, calificacion FROM notas WHERE id = ?`, id)
	var item models.Nota
	err := row.Scan(&item.ID, &item.Estudiante, &item.Curso, &item.Calificacion)
	return item, err
}

func (r NotaRepository) Create(ctx context.Context, item models.Nota) (models.Nota, error) {
	result, err := r.DB.ExecContext(ctx, `INSERT INTO notas (estudiante, curso, calificacion) VALUES (?, ?, ?)`, item.Estudiante, item.Curso, item.Calificacion)
	if err != nil {
		return models.Nota{}, err
	}
	item.ID, err = result.LastInsertId()
	return item, err
}

func (r NotaRepository) Update(ctx context.Context, item models.Nota) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE notas SET estudiante = ?, curso = ?, calificacion = ? WHERE id = ?`, item.Estudiante, item.Curso, item.Calificacion, item.ID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r NotaRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM notas WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}

// GuardarMemento inserta un snapshot inmutable en la tabla de versiones.
func (r NotaRepository) GuardarMemento(memento models.NotaMemento) (models.NotaMemento, error) {
	result, err := r.DB.Exec(`INSERT INTO nota_versiones (nota_id, estudiante, curso, calificacion, fecha_modificacion) VALUES (?, ?, ?, ?, ?)`,
		memento.NotaID, memento.Estudiante, memento.Curso, memento.Calificacion, memento.FechaModificacion)
	if err != nil {
		return models.NotaMemento{}, err
	}
	memento.ID, err = result.LastInsertId()
	return memento, err
}

// ObtenerHistorial retorna todas las versiones anteriores ordenadas de la más reciente a la más antigua.
func (r NotaRepository) ObtenerHistorial(notaID int64) ([]models.NotaMemento, error) {
	rows, err := r.DB.Query(`SELECT id, nota_id, estudiante, curso, calificacion, fecha_modificacion FROM nota_versiones WHERE nota_id = ? ORDER BY id DESC`, notaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.NotaMemento, 0)
	for rows.Next() {
		var m models.NotaMemento
		if err := rows.Scan(&m.ID, &m.NotaID, &m.Estudiante, &m.Curso, &m.Calificacion, &m.FechaModificacion); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

// ObtenerMementoPorID recupera una versión específica para su restauración.
func (r NotaRepository) ObtenerMementoPorID(versionID int64) (models.NotaMemento, error) {
	row := r.DB.QueryRow(`SELECT id, nota_id, estudiante, curso, calificacion, fecha_modificacion FROM nota_versiones WHERE id = ?`, versionID)
	var m models.NotaMemento
	err := row.Scan(&m.ID, &m.NotaID, &m.Estudiante, &m.Curso, &m.Calificacion, &m.FechaModificacion)
	return m, err
}
