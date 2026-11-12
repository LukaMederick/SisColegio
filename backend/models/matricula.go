package models

import "fmt"

type Matricula struct {
	ID         int64  `json:"id"`
	Estudiante string `json:"estudiante"`
	Curso      string `json:"curso"`
	Estado     string `json:"estado"`
}

type Estado interface{ PuedePasar(string) bool }
type Pendiente struct{}
func (Pendiente) PuedePasar(n string) bool { return n == "Activa" || n == "Anulada" }
type Activa struct{}
func (Activa) PuedePasar(n string) bool { return n == "Finalizada" || n == "Anulada" }
type EstadoFinal struct{}
func (EstadoFinal) PuedePasar(string) bool { return false }

func ObtenerEstado(nom string) Estado {
	switch nom {
	case "Activa": return Activa{}
	case "Finalizada", "Anulada": return EstadoFinal{}
	default: return Pendiente{}
	}
}

func (m *Matricula) CambiarEstado(nuevo string) error {
	if !ObtenerEstado(m.Estado).PuedePasar(nuevo) {
		return fmt.Errorf("transición inválida de %s a %s", m.Estado, nuevo)
	}
	m.Estado = nuevo
	return nil
}
