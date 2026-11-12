package main

// Patrón STATE: Control del ciclo de vida de una matrícula
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

type Final struct{ nombre string }
func (f Final) Nombre() string { return f.nombre }
func (Final) PuedePasar(string) bool { return false }

func ObtenerEstado(nom string) Estado {
	switch nom {
	case "Activa": return Activa{}
	case "Finalizada": return Final{"Finalizada"}
	case "Anulada": return Final{"Anulada"}
	default: return Pendiente{}
	}
}

