package entities

// DoctorEvent is the domain event fired by the doctor feature.
// Inject an events.EventPublisher to broadcast it via WebSocket.
//
// Event types:
//   - "estado_actualizado"  → doctor changed a patient's clinical state
//   - "nota_clinica"        → doctor added a new PatientNote
type DoctorEvent struct {
	// Tipo identifies what happened so the WebSocket consumer can route it correctly.
	Tipo string `json:"tipo"`

	// Patient context
	PacienteID     string `json:"paciente_id"`
	NombrePaciente string `json:"nombre_paciente"`

	// Who triggered the event
	DoctorID     string `json:"doctor_id"`
	NombreDoctor string `json:"nombre_doctor"`

	// Estado fields (populated for "estado_actualizado")
	NuevoEstado         string `json:"nuevo_estado,omitempty"`
	UltimaActualizacion string `json:"ultima_actualizacion,omitempty"`

	// Note fields (populated for "nota_clinica")
	// The full note is embedded so the WebSocket consumer can render it without
	// a secondary DB query.
	Nota *PatientNote `json:"nota,omitempty"`
}
