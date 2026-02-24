package application

import (
	"ApiMedicGO/src/feature/pacientes/domain/entities"
	"ApiMedicGO/src/feature/pacientes/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePacienteUseCase struct {
	repo repositories.PacienteRepository
}

func NewCreatePacienteUseCase(repo repositories.PacienteRepository) *CreatePacienteUseCase {
	return &CreatePacienteUseCase{repo: repo}
}

// CreatePacienteRequest DTO de entrada para crear un paciente.
type CreatePacienteRequest struct {
	Nombre          string  `json:"nombre"           binding:"required"`
	Edad            int     `json:"edad"             binding:"required,min=0,max=150"`
	Estatura        float64 `json:"estatura"         binding:"required,min=1"`   // cm
	Peso            float64 `json:"peso"             binding:"required,min=0.1"` // kg
	RitmoCardiaco   int     `json:"ritmo_cardiaco"   binding:"required,min=0"`   // BPM
	Saturacion      float64 `json:"saturacion"       binding:"required,min=0,max=100"`
	PresionArterial string  `json:"presion_arterial" binding:"required"` // ej. "120/80"
	SpO2            float64 `json:"spo2"             binding:"required,min=0,max=100"`
}

func (uc *CreatePacienteUseCase) Execute(req CreatePacienteRequest) (*entities.Paciente, error) {
	paciente := &entities.Paciente{
		ID:              primitive.NewObjectID(),
		Nombre:          req.Nombre,
		Edad:            req.Edad,
		Estatura:        req.Estatura,
		Peso:            req.Peso,
		RitmoCardiaco:   req.RitmoCardiaco,
		Saturacion:      req.Saturacion,
		PresionArterial: req.PresionArterial,
		SpO2:            req.SpO2,
	}

	if err := uc.repo.Save(paciente); err != nil {
		return nil, err
	}
	return paciente, nil
}
