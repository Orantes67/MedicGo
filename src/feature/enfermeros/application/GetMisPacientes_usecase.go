package application

import (
	"ApiMedicGO/src/feature/enfermeros/domain/entities"
	"ApiMedicGO/src/feature/enfermeros/domain/repositories"
)

// GetMisPacientesUseCase returns all patients assigned to the authenticated nurse
// together with the summary stats shown on the Home screen.
type GetMisPacientesUseCase struct {
	repo repositories.EnfermeroRepository
}

func NewGetMisPacientesUseCase(repo repositories.EnfermeroRepository) *GetMisPacientesUseCase {
	return &GetMisPacientesUseCase{repo: repo}
}

// Execute fetches the nurse's patients and computes the stats counters.
func (uc *GetMisPacientesUseCase) Execute(enfermeroID string) (*entities.MisPacientesResponse, error) {
	pacientes, err := uc.repo.GetPacientesAsignados(enfermeroID)
	if err != nil {
		return nil, err
	}

	stats := entities.ResumenStats{Total: len(pacientes)}
	for _, p := range pacientes {
		switch p.EstadoActual {
		case entities.EstadoCritico:
			stats.Criticos++
		case entities.EstadoObservacion:
			stats.Observacion++
		case entities.EstadoEstable:
			stats.Estables++
		}
	}

	return &entities.MisPacientesResponse{
		Stats:     stats,
		Pacientes: pacientes,
	}, nil
}
