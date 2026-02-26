package application

import (
	"errors"
	"fmt"
	"time"

	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/feature/doctores/domain/entities"
	"ApiMedicGO/src/feature/doctores/domain/repositories"
	enfEntities "ApiMedicGO/src/feature/enfermeros/domain/entities"
)

// UpdateEstadoPacienteDoctorUseCase lets a doctor change a patient's clinical state.
// After persisting it fires a DoctorEvent so the WebSocket hub can broadcast in real time.
type UpdateEstadoPacienteDoctorUseCase struct {
	repo      repositories.DoctorRepository
	publisher events.EventPublisher
}

func NewUpdateEstadoPacienteDoctorUseCase(
	repo repositories.DoctorRepository,
	publisher events.EventPublisher,
) *UpdateEstadoPacienteDoctorUseCase {
	return &UpdateEstadoPacienteDoctorUseCase{repo: repo, publisher: publisher}
}

// UpdateEstadoDoctorRequest is the DTO sent by the doctor from the Cards screen.
type UpdateEstadoDoctorRequest struct {
	NuevoEstado string `json:"nuevo_estado" binding:"required"`
}

// Execute validates, persists, and publishes the "estado_actualizado" event.
func (uc *UpdateEstadoPacienteDoctorUseCase) Execute(
	pacienteID string,
	doctorID string,
	nombreDoctor string,
	nombrePaciente string,
	req UpdateEstadoDoctorRequest,
) error {
	switch req.NuevoEstado {
	case enfEntities.EstadoEstable, enfEntities.EstadoCritico, enfEntities.EstadoObservacion:
		// valid
	default:
		return fmt.Errorf("estado inválido: debe ser '%s', '%s' o '%s'",
			enfEntities.EstadoEstable, enfEntities.EstadoCritico, enfEntities.EstadoObservacion)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err := uc.repo.UpdateEstadoPaciente(pacienteID, doctorID, req.NuevoEstado, timestamp); err != nil {
		return errors.New("no se pudo actualizar el estado: " + err.Error())
	}

	_ = uc.publisher.Publish(entities.DoctorEvent{
		Tipo:                "estado_actualizado",
		PacienteID:          pacienteID,
		NombrePaciente:      nombrePaciente,
		DoctorID:            doctorID,
		NombreDoctor:        nombreDoctor,
		NuevoEstado:         req.NuevoEstado,
		UltimaActualizacion: timestamp,
	})

	return nil
}
