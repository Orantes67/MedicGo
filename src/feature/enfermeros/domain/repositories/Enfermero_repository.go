package repositories

import (
	"ApiMedicGO/src/feature/enfermeros/domain/entities"
	doctorEntities "ApiMedicGO/src/feature/doctores/domain/entities"
)

// EnfermeroRepository defines the data operations needed by the enfermero feature.
// The underlying store queries the pacientes and notas_clinicas collections.
type EnfermeroRepository interface {
	// GetPacientesAsignados returns all patients currently assigned to the given nurse.
	GetPacientesAsignados(enfermeroID string) ([]*entities.PacienteResumen, error)

	// UpdateEstadoPaciente atomically updates a patient's state, condition note, and
	// last-update timestamp.  enfermeroID is used as an extra filter so a nurse can
	// only modify patients that are actually assigned to her.
	UpdateEstadoPaciente(pacienteID string, enfermeroID string, nuevoEstado string, nota string, timestamp string) error

	// GetNotasByPaciente returns all clinical notes stored for a patient.
	// enfermeroID is used to verify the patient is assigned to this nurse.
	GetNotasByPaciente(pacienteID string, enfermeroID string) ([]*doctorEntities.PatientNote, error)

	// AddNota persists a new clinical note in the notas_clinicas collection.
	// The PatientNote.DoctorID field carries the enfermero's ObjectID as the author.
	AddNota(nota *doctorEntities.PatientNote) (*doctorEntities.PatientNote, error)
}
