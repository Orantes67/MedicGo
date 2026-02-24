package application

import "ApiMedicGO/src/feature/pacientes/domain/repositories"

type DeletePacienteUseCase struct {
	repo repositories.PacienteRepository
}

func NewDeletePacienteUseCase(repo repositories.PacienteRepository) *DeletePacienteUseCase {
	return &DeletePacienteUseCase{repo: repo}
}

func (uc *DeletePacienteUseCase) Execute(id string) error {
	return uc.repo.Delete(id)
}
