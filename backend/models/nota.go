package models

import "time"

// NotaMemento representa el Memento en el patrón.
// Es una copia inmutable del estado de una Nota con su marca temporal.
type NotaMemento struct {
	ID                int64   `json:"id"`
	NotaID            int64   `json:"nota_id"`
	Estudiante        string  `json:"estudiante"`
	Curso             string  `json:"curso"`
	Calificacion      float64 `json:"calificacion"`
	FechaModificacion string  `json:"fecha_modificacion"`
}

// Nota representa el Originator en el patrón Memento.
// Tiene la responsabilidad de crear snapshots de su estado (mementos) y restaurarse a partir de ellos.
type Nota struct {
	ID           int64   `json:"id"`
	Estudiante   string  `json:"estudiante"`
	Curso        string  `json:"curso"`
	Calificacion float64 `json:"calificacion"`
}

// CreateMemento genera un nuevo Memento con el estado actual de la Nota.
func (n Nota) CreateMemento() NotaMemento {
	return NotaMemento{
		NotaID:            n.ID,
		Estudiante:        n.Estudiante,
		Curso:             n.Curso,
		Calificacion:      n.Calificacion,
		FechaModificacion: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// Restore recupera el estado del Originator desde un Memento sin alterar su ID principal.
func (n *Nota) Restore(m NotaMemento) {
	n.Estudiante = m.Estudiante
	n.Curso = m.Curso
	n.Calificacion = m.Calificacion
}

// CaretakerStore define las operaciones necesarias para la persistencia de Mementos.
type CaretakerStore interface {
	GuardarMemento(memento NotaMemento) (NotaMemento, error)
	ObtenerHistorial(notaID int64) ([]NotaMemento, error)
	ObtenerMementoPorID(versionID int64) (NotaMemento, error)
}

// NotaCaretaker administra el historial de versiones (Mementos).
// No modifica ni interpreta la lógica interna del Originator ni del Memento.
type NotaCaretaker struct {
	Store CaretakerStore
}

// NuevoNotaCaretaker crea una nueva instancia del Caretaker.
func NuevoNotaCaretaker(store CaretakerStore) *NotaCaretaker {
	return &NotaCaretaker{Store: store}
}

// GuardarSnapshot solicita almacenar un Memento en el historial.
func (c *NotaCaretaker) GuardarSnapshot(memento NotaMemento) (NotaMemento, error) {
	return c.Store.GuardarMemento(memento)
}

// ObtenerHistorial recupera la lista cronológica de Mementos de una nota.
func (c *NotaCaretaker) ObtenerHistorial(notaID int64) ([]NotaMemento, error) {
	return c.Store.ObtenerHistorial(notaID)
}

// ObtenerVersion recupera un Memento específico por su identificador.
func (c *NotaCaretaker) ObtenerVersion(versionID int64) (NotaMemento, error) {
	return c.Store.ObtenerMementoPorID(versionID)
}
