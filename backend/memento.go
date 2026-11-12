package main

// Patrón MEMENTO: Originator y Memento para restaurar calificaciones
type Nota struct {
	ID           int64   `json:"id"`
	Estudiante   string  `json:"estudiante"`
	Curso        string  `json:"curso"`
	Calificacion float64 `json:"calificacion"`
}

type NotaMemento struct {
	NotaID       int64   `json:"nota_id"`
	Calificacion float64 `json:"calificacion"`
}

func (n Nota) CrearMemento() NotaMemento {
	return NotaMemento{NotaID: n.ID, Calificacion: n.Calificacion}
}

func (n *Nota) Restaurar(m NotaMemento) {
	n.Calificacion = m.Calificacion
}

