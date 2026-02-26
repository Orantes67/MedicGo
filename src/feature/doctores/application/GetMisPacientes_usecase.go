package application

import (
	"ApiMedicGO/src/feature/doctores/domain/entities"
	"ApiMedicGO/src/feature/doctores/domain/repositories"
	enfEntities "ApiMedicGO/src/feature/enfermeros/domain/entities"
)

// GetMisPacientesDoctorUseCase returns all patients assigned to the authenticated doctor,
// the top-panel stats, and the "Prioridad" slice (critical patients only).
type GetMisPacientesDoctorUseCase struct {
	repo repositories.DoctorRepository
}

func NewGetMisPacientesDoctorUseCase(repo repositories.DoctorRepository) *GetMisPacientesDoctorUseCase {
	return &GetMisPacientesDoctorUseCase{repo: repo}
}

// Execute fetches the list, computes stats, and separates the priority (critical) patients.
func (uc *GetMisPacientesDoctorUseCase) Execute(doctorID string) (*entities.MisPacientesDoctorResponse, error) {
	pacientes, err := uc.repo.GetPacientesAsignados(doctorID)
	if err != nil {
		return nil, err
	}

	stats := entities.ResumenStatsDoctor{Total: len(pacientes)}
	var prioridad []*entities.PacienteResumenDoctor

	for _, p := range pacientes {
		switch p.EstadoActual {
		case enfEntities.EstadoCritico:
			stats.Criticos++
			prioridad = append(prioridad, p) // red-highlight section
		case enfEntities.EstadoObservacion:
			stats.Observacion++
		case enfEntities.EstadoEstable:
			stats.Estables++
		}
	}

	return &entities.MisPacientesDoctorResponse{
		Stats:     stats,
		Prioridad: prioridad,
		Pacientes: pacientes,
	}, nil
}
