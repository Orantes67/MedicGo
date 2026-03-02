package entities

// PacienteResumen is the read-model used in the nurse's patient list (Home screen).
// It contains only the fields the nurse needs to see per card.
type PacienteResumen struct {
	ID                  string `json:"id"`
	Nombre              string `json:"nombre"`
	Apellido            string `json:"apellido"`
	Edad                int    `json:"edad"`
	AreaNombre          string `json:"area_nombre"`
	EstadoActual        string `json:"estado_actual"`
	NotaCondicion       string `json:"nota_condicion"`
	NombreDoctor        string `json:"nombre_doctor"`
	UltimaActualizacion string `json:"ultima_actualizacion"`
}

// ResumenStats holds the summary counters shown at the top of the Home screen.
type ResumenStats struct {
	Total      int `json:"total"`
	Criticos   int `json:"criticos"`
	Observacion int `json:"observacion"`
	Estables   int `json:"estables"`
}

// MisPacientesResponse is the full payload returned by GET /enfermeros/mis-pacientes.
type MisPacientesResponse struct {
	Stats    ResumenStats      `json:"stats"`
	Pacientes []*PacienteResumen `json:"pacientes"`
}
