package models

type Docente struct {
	ID           int64  `json:"id"`
	Nombre       string `json:"nombre"`
	Especialidad string `json:"especialidad"`
	Correo       string `json:"correo"`
}

type Clase struct {
	ID            int64  `json:"id"`
	Nombre        string `json:"nombre"`
	Horario       string `json:"horario"`
	DocenteID     int64  `json:"docenteId"`
	DocenteNombre string `json:"docenteNombre,omitempty"`
}

type Asistencia struct {
	ID          int64  `json:"id"`
	MatriculaID int64  `json:"matriculaId"`
	Estudiante  string `json:"estudiante,omitempty"`
	ClaseID     int64  `json:"claseId"`
	Clase       string `json:"clase,omitempty"`
	Fecha       string `json:"fecha"`
	Estado      string `json:"estado"`
}

var EstadosAsistencia = map[string]bool{
	"Presente": true,
	"Tardanza": true,
	"Ausente":  true,
}
