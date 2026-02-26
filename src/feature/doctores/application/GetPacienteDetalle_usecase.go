package application

import (
	"ApiMedicGO/src/feature/doctores/domain/entities"
	"ApiMedicGO/src/feature/doctores/domain/repositories"
)

// GetPacienteDetalleUseCase returns the rich patient view for the Cards screen,
// including all clinical notes the doctor has written for this patient.
type GetPacienteDetalleUseCase struct {
	repo repositories.DoctorRepository
}

func NewGetPacienteDetalleUseCase(repo repositories.DoctorRepository) *GetPacienteDetalleUseCase {
	return &GetPacienteDetalleUseCase{repo: repo}
}

// Execute returns the full patient model.  Returns an error if the patient is
// not assigned to this doctor (enforced at the DB query level).
func (uc *GetPacienteDetalleUseCase) Execute(pacienteID string, doctorID string) (*entities.PacienteDetalleDoctor, error) {
	return uc.repo.GetPacienteDetalle(pacienteID, doctorID)
}
