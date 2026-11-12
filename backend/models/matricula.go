package models

import "fmt"

// EstadoMatricula define la interfaz del patrón State para el ciclo de vida de una matrícula.
type EstadoMatricula interface {
	Nombre() string
	PuedeCambiarA(destino string) bool
	Acciones() []string
	Transicionar(destino string) (EstadoMatricula, error)
}

// EstadoActiva representa la matrícula de un estudiante regular activo.
type EstadoActiva struct{}

func (e EstadoActiva) Nombre() string { return "Activa" }
func (e EstadoActiva) PuedeCambiarA(destino string) bool {
	return destino == "Finalizada" || destino == "Anulada"
}
func (e EstadoActiva) Acciones() []string { return []string{"Finalizada", "Anulada"} }
func (e EstadoActiva) Transicionar(destino string) (EstadoMatricula, error) {
	if !e.PuedeCambiarA(destino) {
		return nil, fmt.Errorf("transición inválida: de Activa a %s", destino)
	}
	return ObtenerEstadoMatricula(destino)
}

// EstadoFinalizada representa el estado terminal al concluir el año lectivo.
type EstadoFinalizada struct{}

func (e EstadoFinalizada) Nombre() string { return "Finalizada" }
func (e EstadoFinalizada) PuedeCambiarA(destino string) bool { return false }
func (e EstadoFinalizada) Acciones() []string { return []string{} }
func (e EstadoFinalizada) Transicionar(destino string) (EstadoMatricula, error) {
	return nil, fmt.Errorf("la matrícula ya está Finalizada y no permite cambios")
}

// EstadoAnulada representa el estado terminal por retiro, traslado o cancelación.
type EstadoAnulada struct{}

func (e EstadoAnulada) Nombre() string { return "Anulada" }
func (e EstadoAnulada) PuedeCambiarA(destino string) bool { return false }
func (e EstadoAnulada) Acciones() []string { return []string{} }
func (e EstadoAnulada) Transicionar(destino string) (EstadoMatricula, error) {
	return nil, fmt.Errorf("la matrícula está Anulada y no permite cambios")
}

// EstadoPendiente representa el estado previo a la confirmación de la matrícula.
type EstadoPendiente struct{}

func (e EstadoPendiente) Nombre() string { return "Pendiente" }
func (e EstadoPendiente) PuedeCambiarA(destino string) bool {
	return destino == "Activa" || destino == "Anulada"
}
func (e EstadoPendiente) Acciones() []string { return []string{"Activa", "Anulada"} }
func (e EstadoPendiente) Transicionar(destino string) (EstadoMatricula, error) {
	if !e.PuedeCambiarA(destino) {
		return nil, fmt.Errorf("transición inválida: de Pendiente a %s", destino)
	}
	return ObtenerEstadoMatricula(destino)
}

// EstadosMatricula lista los estados reconocidos por el sistema.
var EstadosMatricula = map[string]bool{
	"Pendiente":  true,
	"Activa":     true,
	"Finalizada": true,
	"Anulada":    true,
}

// ObtenerEstadoMatricula es la factoría que retorna la clase de estado correspondiente.
func ObtenerEstadoMatricula(nombre string) (EstadoMatricula, error) {
	switch nombre {
	case "Activa":
		return EstadoActiva{}, nil
	case "Finalizada":
		return EstadoFinalizada{}, nil
	case "Anulada":
		return EstadoAnulada{}, nil
	case "Pendiente":
		return EstadoPendiente{}, nil
	default:
		return nil, fmt.Errorf("estado no reconocido: %s", nombre)
	}
}

// Matricula representa la inscripción de un estudiante.
type Matricula struct {
	ID                int64  `json:"id"`
	Estudiante        string `json:"estudiante"`
	Apellido          string `json:"apellido"`
	GradoCurso        string `json:"grado_curso"`
	Seccion           string `json:"seccion"`
	FechaInscripcion  string `json:"fecha_inscripcion"`
	TutorResponsable  string `json:"tutor_responsable"`
	Estado            string `json:"estado"`
}

// CambiarEstado aplica el patrón State delegando la transición en el objeto de estado actual.
func (m *Matricula) CambiarEstado(nuevoEstado string) error {
	estadoActual, err := ObtenerEstadoMatricula(m.Estado)
	if err != nil {
		return err
	}
	siguiente, err := estadoActual.Transicionar(nuevoEstado)
	if err != nil {
		return err
	}
	m.Estado = siguiente.Nombre()
	return nil
}

// PuedeCambiarEstadoMatricula verifica si la transición es admisible según la clase de estado.
func PuedeCambiarEstadoMatricula(actual, siguiente string) bool {
	estado, err := ObtenerEstadoMatricula(actual)
	if err != nil {
		return false
	}
	return estado.PuedeCambiarA(siguiente)
}
