package application

import (
	"ApiMedicGO/src/feature/admin/domain/entities"
	"ApiMedicGO/src/feature/admin/domain/repositories"
)

// GetMetricsUseCase devuelve las métricas del panel principal del administrador.
type GetMetricsUseCase struct {
	repo repositories.AdminRepository
}

func NewGetMetricsUseCase(repo repositories.AdminRepository) *GetMetricsUseCase {
	return &GetMetricsUseCase{repo: repo}
}

func (uc *GetMetricsUseCase) Execute() (*entities.AdminMetrics, error) {
	return uc.repo.GetMetrics()
}
