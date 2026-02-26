package entities

// EstadoPaciente groups the valid patient state values.
const (
	EstadoEstable     = "Estable"
	EstadoCritico     = "Crítico"
	EstadoObservacion = "Observación"
)

// PacienteEstadoEvent is the domain event emitted every time a nurse
// updates a patient's state or adds a clinical / quick note.
//
// Connect this to your WebSocket hub by implementing the
// events.EventPublisher interface in your external project.
//
// Event types:
//   - "estado_actualizado" → nurse changed the patient state
//   - "nota_rapida"        → nurse added a quick observation note (UpdateEstado flow)
//   - "nota_clinica"       → nurse added a structured clinical note (AddNota flow)
type PacienteEstadoEvent struct {
	// Tipo distinguishes the kind of event so consumers can react accordingly.
	Tipo string `json:"tipo"`

	// Patient context
	PacienteID     string `json:"paciente_id"`
	NombrePaciente string `json:"nombre_paciente"`

	// Who triggered the event
	EnfermeroID     string `json:"enfermero_id"`
	NombreEnfermero string `json:"nombre_enfermero"`

	// Estado fields (populated for "estado_actualizado" and "nota_rapida")
	NuevoEstado         string `json:"nuevo_estado,omitempty"`
	NotaRapida          string `json:"nota_rapida,omitempty"`
	UltimaActualizacion string `json:"ultima_actualizacion,omitempty"`

	// Note fields (populated for "nota_clinica").
	// Inline fields avoid a cross-domain import in the entity layer while
	// giving the WebSocket consumer enough data to render the note immediately.
	NoteID      string `json:"note_id,omitempty"`
	NoteContent string `json:"note_content,omitempty"`
	NoteDate    string `json:"note_date,omitempty"`
}
