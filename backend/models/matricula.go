package models

import "fmt"

// Entidad Matrícula
type Matricula struct {
	ID         int64  `json:"id"`
	Estudiante string `json:"estudiante"`
	Curso      string `json:"curso"`
	Estado     string `json:"estado"`
}

// PATRÓN STATE: Interfaz que desacopla la lógica de transición
type Estado interface{ PuedePasar(string) bool }

// PATRÓN STATE: Estado Pendiente -> solo puede ir a Activa o Anulada
type Pendiente struct{}
func (Pendiente) PuedePasar(n string) bool { return n == "Activa" || n == "Anulada" }

// PATRÓN STATE: Estado Activa -> solo puede ir a Finalizada o Anulada
type Activa struct{}
func (Activa) PuedePasar(n string) bool { return n == "Finalizada" || n == "Anulada" }

// PATRÓN STATE: Estados terminales (no permiten más cambios)
type EstadoFinal struct{}
func (EstadoFinal) PuedePasar(string) bool { return false }

// Factoría para instanciar el objeto State según el nombre actual
func ObtenerEstado(nom string) Estado {
	switch nom {
	case "Activa": return Activa{}
	case "Finalizada", "Anulada": return EstadoFinal{}
	default: return Pendiente{}
	}
}

// Delegación al State: si el estado actual no lo permite, rechaza la operación
func (m *Matricula) CambiarEstado(nuevo string) error {
	if !ObtenerEstado(m.Estado).PuedePasar(nuevo) {
		return fmt.Errorf("transición inválida de %s a %s", m.Estado, nuevo)
	}
	m.Estado = nuevo
	return nil
}
