package repositories

import "ApiMedicGO/src/feature/doctores/domain/entities"

// DoctorRepository defines all data operations needed by the doctor feature.
type DoctorRepository interface {
	// GetPacientesAsignados returns all patients assigned to the given doctor.
	GetPacientesAsignados(doctorID string) ([]*entities.PacienteResumenDoctor, error)

	// GetPacienteDetalle returns the full read-model for the Cards screen,
	// including all clinical notes written by this doctor for the patient.
	// Returns an error if the patient is not assigned to this doctor.
	GetPacienteDetalle(pacienteID string, doctorID string) (*entities.PacienteDetalleDoctor, error)

	// UpdateEstadoPaciente atomically changes the patient's clinical state and
	// last-update timestamp.  doctorID is included in the filter so a doctor
	// can only update patients that belong to him.
	UpdateEstadoPaciente(pacienteID string, doctorID string, nuevoEstado string, timestamp string) error

	// AddNota persists a new clinical note in the "notas_clinicas" collection.
	AddNota(nota *entities.PatientNote) (*entities.PatientNote, error)

	// GetNotasByPaciente returns all notes that this doctor wrote for the given patient.
	GetNotasByPaciente(pacienteID string, doctorID string) ([]*entities.PatientNote, error)
}
