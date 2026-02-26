package entities

// PacienteResumenDoctor is the read-model for each card in the doctor's Home screen.
// It includes the nurse name so the detail card shows "Enfermera Asignada" immediately.
type PacienteResumenDoctor struct {
	ID                  string `json:"id"`
	Nombre              string `json:"nombre"`
	Apellido            string `json:"apellido"`
	Edad                int    `json:"edad"`
	AreaNombre          string `json:"area_nombre"`
	EstadoActual        string `json:"estado_actual"`
	NotaCondicion       string `json:"nota_condicion"`
	NombreEnfermero     string `json:"nombre_enfermero"`
	UltimaActualizacion string `json:"ultima_actualizacion"`
	FechaRegistro       string `json:"fecha_registro"`
}

// ResumenStatsDoctor holds the counters shown in the top summary panel.
type ResumenStatsDoctor struct {
	Total       int `json:"total"`
	Criticos    int `json:"criticos"`
	Observacion int `json:"observacion"`
	Estables    int `json:"estables"`
}

// MisPacientesDoctorResponse is the full payload for GET /doctores/mis-pacientes.
// Prioridad lists only the critical patients so the UI can render the red-highlight section.
type MisPacientesDoctorResponse struct {
	Stats     ResumenStatsDoctor      `json:"stats"`
	Prioridad []*PacienteResumenDoctor `json:"prioridad"`
	Pacientes []*PacienteResumenDoctor `json:"pacientes"`
}

// PacienteDetalleDoctor is the rich read-model for the doctor's Cards screen.
// It includes the full clinical notes list so the UI can render them inline.
type PacienteDetalleDoctor struct {
	ID                  string         `json:"id"`
	Nombre              string         `json:"nombre"`
	Apellido            string         `json:"apellido"`
	Edad                int            `json:"edad"`
	AreaNombre          string         `json:"area_nombre"`
	EstadoActual        string         `json:"estado_actual"`
	NombreEnfermero     string         `json:"nombre_enfermero"`
	UltimaActualizacion string         `json:"ultima_actualizacion"`
	NotasClinicas       []*PatientNote `json:"notas_clinicas"`
}
