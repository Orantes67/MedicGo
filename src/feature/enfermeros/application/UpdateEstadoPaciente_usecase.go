package application

import (
	"errors"
	"fmt"
	"time"

	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/feature/enfermeros/domain/entities"
	"ApiMedicGO/src/feature/enfermeros/domain/repositories"
)

// UpdateEstadoPacienteUseCase lets a nurse change a patient's state and leave
// a quick observation note.  After persisting the change it fires a
// PacienteEstadoEvent so the WebSocket hub (external project) can broadcast
// the update in real time.
type UpdateEstadoPacienteUseCase struct {
	repo      repositories.EnfermeroRepository
	publisher events.EventPublisher
}

func NewUpdateEstadoPacienteUseCase(
	repo repositories.EnfermeroRepository,
	publisher events.EventPublisher,
) *UpdateEstadoPacienteUseCase {
	return &UpdateEstadoPacienteUseCase{repo: repo, publisher: publisher}
}

// UpdateEstadoRequest is the DTO sent by the nurse from the Cards screen.
type UpdateEstadoRequest struct {
	NuevoEstado string `json:"nuevo_estado" binding:"required"`
	NotaRapida  string `json:"nota_rapida"`
}

// Execute validates, persists and publishes the state-change event.
//
// Parameters:
//   - pacienteID     → MongoDB ObjectID of the patient (path param)
//   - enfermeroID    → user_id from the JWT (ensures ownership)
//   - nombreEnfermero → license / display name from the JWT
//   - nombrePaciente  → denormalized name for the event payload
func (uc *UpdateEstadoPacienteUseCase) Execute(
	pacienteID string,
	enfermeroID string,
	nombreEnfermero string,
	nombrePaciente string,
	req UpdateEstadoRequest,
) error {
	// Validate state value
	switch req.NuevoEstado {
	case entities.EstadoEstable, entities.EstadoCritico, entities.EstadoObservacion:
		// valid
	default:
		return fmt.Errorf("estado inválido: debe ser '%s', '%s' o '%s'",
			entities.EstadoEstable, entities.EstadoCritico, entities.EstadoObservacion)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err := uc.repo.UpdateEstadoPaciente(
		pacienteID, enfermeroID, req.NuevoEstado, req.NotaRapida, timestamp,
	); err != nil {
		return errors.New("no se pudo actualizar el estado: " + err.Error())
	}

	// ─── Domain event ────────────────────────────────────────────────────────
	// Determine event type: if a quick note was provided it counts as "nota_rapida",
	// otherwise it is a plain "estado_actualizado".
	tipo := "estado_actualizado"
	if req.NotaRapida != "" {
		tipo = "nota_rapida"
	}

	event := entities.PacienteEstadoEvent{
		Tipo:                tipo,
		PacienteID:          pacienteID,
		NombrePaciente:      nombrePaciente,
		EnfermeroID:         enfermeroID,
		NombreEnfermero:     nombreEnfermero,
		NuevoEstado:         req.NuevoEstado,
		NotaRapida:          req.NotaRapida,
		UltimaActualizacion: timestamp,
	}

	// Publish is fire-and-forget; a failure does NOT roll back the DB update.
	// The WebSocket consumer can retry or log independently.
	_ = uc.publisher.Publish(event)

	return nil
}
