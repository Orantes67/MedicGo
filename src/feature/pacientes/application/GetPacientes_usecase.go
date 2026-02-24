package application

import (
	"errors"

	"ApiMedicGO/src/feature/pacientes/domain/entities"
	"ApiMedicGO/src/feature/pacientes/domain/repositories"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
)

type GetPacientesUseCase struct {
	repo repositories.PacienteRepository
}

func NewGetPacientesUseCase(repo repositories.PacienteRepository) *GetPacientesUseCase {
	return &GetPacientesUseCase{repo: repo}
}

// Execute retorna pacientes según el rol del usuario:
//   - jefe_doctor / jefe_enfermera → todos los pacientes
//   - doctor   → solo pacientes asignados a él
//   - enfermero → solo pacientes asignados a él
func (uc *GetPacientesUseCase) Execute(role, userID string) ([]*entities.Paciente, error) {
	switch role {
	case loginEntities.RoleJefeDoctor, loginEntities.RoleJefeEnfermera:
		return uc.repo.FindAll()
	case loginEntities.RoleDoctor:
		return uc.repo.FindByDoctor(userID)
	case loginEntities.RoleNurse:
		return uc.repo.FindByEnfermero(userID)
	default:
		return nil, errors.New("rol no reconocido")
	}
}
