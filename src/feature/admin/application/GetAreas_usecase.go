package application

import (
	"ApiMedicGO/src/feature/admin/domain/entities"
	"ApiMedicGO/src/feature/admin/domain/repositories"
)

// GetAreasUseCase devuelve la distribución de pacientes por área hospitalaria.
type GetAreasUseCase struct {
	repo repositories.AdminRepository
}

func NewGetAreasUseCase(repo repositories.AdminRepository) *GetAreasUseCase {
	return &GetAreasUseCase{repo: repo}
}

func (uc *GetAreasUseCase) Execute() ([]*entities.AreaDistribucion, error) {
	return uc.repo.GetDistribucionAreas()
}
