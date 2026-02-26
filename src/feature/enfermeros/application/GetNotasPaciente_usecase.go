package application

import (
	"ApiMedicGO/src/feature/enfermeros/domain/repositories"
	doctorEntities "ApiMedicGO/src/feature/doctores/domain/entities"
)

// GetNotasPacienteUseCase returns all clinical notes for a patient assigned to the nurse.
// The nurse sees all notes (written by any doctor or by herself) so she can follow the
// full clinical history from her view.
type GetNotasPacienteUseCase struct {
	repo repositories.EnfermeroRepository
}

func NewGetNotasPacienteUseCase(repo repositories.EnfermeroRepository) *GetNotasPacienteUseCase {
	return &GetNotasPacienteUseCase{repo: repo}
}

// Execute returns the notes list.  The repo verifies the patient is assigned to this
// nurse before querying the notas_clinicas collection.
func (uc *GetNotasPacienteUseCase) Execute(pacienteID string, enfermeroID string) ([]*doctorEntities.PatientNote, error) {
	return uc.repo.GetNotasByPaciente(pacienteID, enfermeroID)
}
