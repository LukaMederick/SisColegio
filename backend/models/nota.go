package models

// Patrón MEMENTO: Originator y Memento para respaldo de calificaciones
type Nota struct {
	ID           int64   `json:"id"`
	Estudiante   string  `json:"estudiante"`
	Curso        string  `json:"curso"`
	Calificacion float64 `json:"calificacion"`
}

type NotaMemento struct {
	ID           int64   `json:"id,omitempty"`
	NotaID       int64   `json:"nota_id"`
	Calificacion float64 `json:"calificacion"`
}

// CrearMemento genera un snapshot del estado actual
func (n Nota) CrearMemento() NotaMemento {
	return NotaMemento{NotaID: n.ID, Calificacion: n.Calificacion}
}

// Restaurar recupera el estado desde un Memento previo
func (n *Nota) Restaurar(m NotaMemento) {
	n.Calificacion = m.Calificacion
}

