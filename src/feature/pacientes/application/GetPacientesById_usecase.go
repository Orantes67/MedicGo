package application

import (
	"ApiMedicGO/src/feature/pacientes/domain/entities"
	"ApiMedicGO/src/feature/pacientes/domain/repositories"
)

type GetPacienteByIDUseCase struct {
	repo repositories.PacienteRepository
}

func NewGetPacienteByIDUseCase(repo repositories.PacienteRepository) *GetPacienteByIDUseCase {
	return &GetPacienteByIDUseCase{repo: repo}
}

func (uc *GetPacienteByIDUseCase) Execute(id string) (*entities.Paciente, error) {
	return uc.repo.FindByID(id)
}
