package application

import (
	"errors"

	"ApiMedicGO/src/feature/pacientes/domain/repositories"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
)

type AssignPacienteUseCase struct {
	repo repositories.PacienteRepository
}

func NewAssignPacienteUseCase(repo repositories.PacienteRepository) *AssignPacienteUseCase {
	return &AssignPacienteUseCase{repo: repo}
}

// AssignPacienteRequest DTO para asignar un paciente.
// El rol del jefe determina qué campo se puede modificar:
//   - jefe_doctor  → DoctorID
//   - jefe_enfermera → EnfermeroID
type AssignPacienteRequest struct {
	DoctorID    *string `json:"doctor_id"`
	EnfermeroID *string `json:"enfermero_id"`
}

// Execute asigna un paciente respetando las reglas de rol.
// role: rol del usuario que realiza la operación.
func (uc *AssignPacienteUseCase) Execute(pacienteID string, req AssignPacienteRequest, role string) error {
	switch role {
	case loginEntities.RoleJefeDoctor:
		if req.DoctorID == nil {
			return errors.New("debes proveer doctor_id para asignar")
		}
		return uc.repo.Assign(pacienteID, req.DoctorID, nil)

	case loginEntities.RoleJefeEnfermera:
		if req.EnfermeroID == nil {
			return errors.New("debes proveer enfermero_id para asignar")
		}
		return uc.repo.Assign(pacienteID, nil, req.EnfermeroID)

	default:
		return errors.New("no tienes permisos para asignar pacientes")
	}
}
