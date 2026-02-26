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
// Solo el rol administrador puede modificar estos campos:
//   - DoctorID        → asigna un doctor al paciente
//   - EnfermeroID     → asigna un enfermero al paciente
//   - NombreEnfermero → nombre del enfermero (denormalizado para consultas rápidas)
type AssignPacienteRequest struct {
	DoctorID        *string `json:"doctor_id"`
	EnfermeroID     *string `json:"enfermero_id"`
	NombreEnfermero string  `json:"nombre_enfermero"`
}

// Execute asigna un paciente respetando las reglas de rol.
// role: rol del usuario que realiza la operación.
func (uc *AssignPacienteUseCase) Execute(pacienteID string, req AssignPacienteRequest, role string) error {
	switch role {
	case loginEntities.RoleAdmin:
		if req.DoctorID == nil && req.EnfermeroID == nil {
			return errors.New("debes proveer doctor_id o enfermero_id para asignar")
		}
		return uc.repo.Assign(pacienteID, req.DoctorID, req.EnfermeroID, req.NombreEnfermero)

	default:
		return errors.New("no tienes permisos para asignar pacientes")
	}
}
