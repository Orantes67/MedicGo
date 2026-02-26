package entities

// ActividadItem representa un paciente en la sección "Actividad Reciente" del dashboard.
type ActividadItem struct {
	NombrePaciente  string `json:"nombre_paciente"`
	Area            string `json:"area"`
	NombreEnfermero string `json:"nombre_enfermero"`
	Estado          string `json:"estado"`
}

// AdminMetrics agrupa todas las métricas del panel administrativo (tab Métricas).
type AdminMetrics struct {
	TotalPacientes    int              `json:"total_pacientes"`
	Criticos          int              `json:"criticos"`
	Observacion       int              `json:"observacion"`
	Estables          int              `json:"estables"`
	PersonalActivo    int              `json:"personal_activo"`
	Enfermeros        int              `json:"enfermeros"`
	Doctores          int              `json:"doctores"`
	Alertas           int              `json:"alertas"`
	ActividadReciente []*ActividadItem `json:"actividad_reciente"`
}

// UsuarioResumen es el DTO de lectura de un usuario (doctor o enfermero) para el admin.
type UsuarioResumen struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Especialidad string `json:"especialidad"`
	Rol         string `json:"rol"`
	DoctorID    string `json:"doctor_id,omitempty"`
	DoctorNombre string `json:"doctor_nombre,omitempty"`
}

// AreaDistribucion representa la distribución de pacientes en un área del hospital.
type AreaDistribucion struct {
	Area     string `json:"area"`
	Total    int    `json:"total"`
	Criticos int    `json:"criticos"`
}
