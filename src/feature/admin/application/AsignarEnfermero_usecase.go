package application

import "ApiMedicGO/src/feature/admin/domain/repositories"

// AsignarEnfermeroUseCase vincula un enfermero a un doctor específico.
type AsignarEnfermeroUseCase struct {
	repo repositories.AdminRepository
}

func NewAsignarEnfermeroUseCase(repo repositories.AdminRepository) *AsignarEnfermeroUseCase {
	return &AsignarEnfermeroUseCase{repo: repo}
}

// AsignarEnfermeroRequest es el DTO de entrada.
type AsignarEnfermeroRequest struct {
	EnfermeroID string `json:"enfermero_id" binding:"required"`
	DoctorID    string `json:"doctor_id"    binding:"required"`
}

func (uc *AsignarEnfermeroUseCase) Execute(req AsignarEnfermeroRequest) error {
	return uc.repo.AsignarEnfermeroADoctor(req.EnfermeroID, req.DoctorID)
}
