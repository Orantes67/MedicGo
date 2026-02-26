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
// Solo se actualizan los campos que vengan con valor no-zero.
type UpdatePacienteRequest struct {
	Nombre              string `json:"nombre"`
	Apellido            string `json:"apellido"`
	Edad                int    `json:"edad"`
	AreaNombre          string `json:"area_nombre"`
	EstadoActual        string `json:"estado_actual"`
	NotaCondicion       string `json:"nota_condicion"`
	TipoSangre          string `json:"tipo_sangre"`
	Sintomas            string `json:"sintomas"`
	UltimaActualizacion string `json:"ultima_actualizacion"`
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
	if req.Apellido != "" {
		existing.Apellido = req.Apellido
	}
	if req.Edad > 0 {
		existing.Edad = req.Edad
	}
	if req.AreaNombre != "" {
		existing.AreaNombre = req.AreaNombre
	}
	if req.EstadoActual != "" {
		existing.EstadoActual = req.EstadoActual
	}
	if req.NotaCondicion != "" {
		existing.NotaCondicion = req.NotaCondicion
	}
	if req.TipoSangre != "" {
		existing.TipoSangre = req.TipoSangre
	}
	if req.Sintomas != "" {
		existing.Sintomas = req.Sintomas
	}
	if req.UltimaActualizacion != "" {
		existing.UltimaActualizacion = req.UltimaActualizacion
	}

	if err := uc.repo.Update(id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
