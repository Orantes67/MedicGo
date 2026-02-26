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
	Nombre        string `json:"nombre"         binding:"required"`
	Apellido      string `json:"apellido"       binding:"required"`
	Edad          int    `json:"edad"           binding:"required,min=0,max=150"`
	AreaNombre    string `json:"area_nombre"    binding:"required"`
	EstadoActual  string `json:"estado_actual"  binding:"required"`
	NotaCondicion string `json:"nota_condicion"`
	TipoSangre    string `json:"tipo_sangre"    binding:"required"`
	Sintomas      string `json:"sintomas"`
	FechaRegistro string `json:"fecha_registro" binding:"required"`
}

func (uc *CreatePacienteUseCase) Execute(req CreatePacienteRequest) (*entities.Paciente, error) {
	paciente := &entities.Paciente{
		ID:                  primitive.NewObjectID(),
		Nombre:              req.Nombre,
		Apellido:            req.Apellido,
		Edad:                req.Edad,
		AreaNombre:          req.AreaNombre,
		EstadoActual:        req.EstadoActual,
		NotaCondicion:       req.NotaCondicion,
		TipoSangre:          req.TipoSangre,
		Sintomas:            req.Sintomas,
		FechaRegistro:       req.FechaRegistro,
		UltimaActualizacion: req.FechaRegistro,
	}

	if err := uc.repo.Save(paciente); err != nil {
		return nil, err
	}
	return paciente, nil
}
