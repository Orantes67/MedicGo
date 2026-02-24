package application

import (
	"ApiMedicGO/src/feature/pacientes/domain/entities"
	"ApiMedicGO/src/feature/pacientes/domain/repositories"
)

type UpdatePacienteUseCase struct {
	repo repositories.PacienteRepository
}

func NewUpdatePacienteUseCase(repo repositories.PacienteRepository) *UpdatePacienteUseCase {
	return &UpdatePacienteUseCase{repo: repo}
}

// UpdatePacienteRequest DTO de entrada para actualizar un paciente.
type UpdatePacienteRequest struct {
	Nombre          string  `json:"nombre"`
	Edad            int     `json:"edad"`
	Estatura        float64 `json:"estatura"`
	Peso            float64 `json:"peso"`
	RitmoCardiaco   int     `json:"ritmo_cardiaco"`
	Saturacion      float64 `json:"saturacion"`
	PresionArterial string  `json:"presion_arterial"`
	SpO2            float64 `json:"spo2"`
}

func (uc *UpdatePacienteUseCase) Execute(id string, req UpdatePacienteRequest) (*entities.Paciente, error) {
	existing, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Aplicar solo los campos enviados
	if req.Nombre != "" {
		existing.Nombre = req.Nombre
	}
	if req.Edad > 0 {
		existing.Edad = req.Edad
	}
	if req.Estatura > 0 {
		existing.Estatura = req.Estatura
	}
	if req.Peso > 0 {
		existing.Peso = req.Peso
	}
	if req.RitmoCardiaco > 0 {
		existing.RitmoCardiaco = req.RitmoCardiaco
	}
	if req.Saturacion > 0 {
		existing.Saturacion = req.Saturacion
	}
	if req.PresionArterial != "" {
		existing.PresionArterial = req.PresionArterial
	}
	if req.SpO2 > 0 {
		existing.SpO2 = req.SpO2
	}

	if err := uc.repo.Update(id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
