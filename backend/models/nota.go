package models

// PATRÓN MEMENTO: Originator - Entidad principal que guarda y restaura su estado
type Nota struct {
	ID           int64   `json:"id"`
	Estudiante   string  `json:"estudiante"`
	Curso        string  `json:"curso"`
	Calificacion float64 `json:"calificacion"`
}

// PATRÓN MEMENTO: Memento - Copia inmutable del estado interno de una calificación
type NotaMemento struct {
	ID           int64   `json:"id,omitempty"`
	NotaID       int64   `json:"nota_id"`
	Calificacion float64 `json:"calificacion"`
}

// PATRÓN MEMENTO: El Originator produce un snapshot con su estado actual
func (n Nota) CrearMemento() NotaMemento {
	return NotaMemento{NotaID: n.ID, Calificacion: n.Calificacion}
}

// PATRÓN MEMENTO: El Originator restaura su calificación desde un Memento previo
func (n *Nota) Restaurar(m NotaMemento) {
	n.Calificacion = m.Calificacion
}
