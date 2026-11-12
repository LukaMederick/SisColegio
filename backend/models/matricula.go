package models

import "fmt"

// Matricula representa la entidad de inscripción
type Matricula struct {
	ID         int64  `json:"id"`
	Estudiante string `json:"estudiante"`
	Curso      string `json:"curso"`
	Estado     string `json:"estado"`
}

// Patrón STATE: Define la interfaz y los estados concretos
type Estado interface {
	Nombre() string
	PuedePasar(nuevo string) bool
}

type Pendiente struct{}
func (Pendiente) Nombre() string { return "Pendiente" }
func (Pendiente) PuedePasar(n string) bool { return n == "Activa" || n == "Anulada" }

type Activa struct{}
func (Activa) Nombre() string { return "Activa" }
func (Activa) PuedePasar(n string) bool { return n == "Finalizada" || n == "Anulada" }

type EstadoFinal struct{ nombre string }
func (f EstadoFinal) Nombre() string { return f.nombre }
func (EstadoFinal) PuedePasar(string) bool { return false }

func ObtenerEstado(nom string) Estado {
	switch nom {
	case "Activa": return Activa{}
	case "Finalizada": return EstadoFinal{"Finalizada"}
	case "Anulada": return EstadoFinal{"Anulada"}
	default: return Pendiente{}
	}
}

func (m *Matricula) CambiarEstado(nuevo string) error {
	actual := ObtenerEstado(m.Estado)
	if !actual.PuedePasar(nuevo) {
		return fmt.Errorf("transición inválida de %s a %s", m.Estado, nuevo)
	}
	m.Estado = nuevo
	return nil
}

