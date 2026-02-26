package application

import (
	"errors"
	"time"

	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/feature/doctores/domain/entities"
	"ApiMedicGO/src/feature/doctores/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddNotaClinicaUseCase persists a new PatientNote and fires a "nota_clinica" DoctorEvent.
type AddNotaClinicaUseCase struct {
	repo      repositories.DoctorRepository
	publisher events.EventPublisher
}

func NewAddNotaClinicaUseCase(
	repo repositories.DoctorRepository,
	publisher events.EventPublisher,
) *AddNotaClinicaUseCase {
	return &AddNotaClinicaUseCase{repo: repo, publisher: publisher}
}

// AddNotaRequest is the DTO sent by the doctor when writing a clinical note.
type AddNotaRequest struct {
	Content string `json:"content" binding:"required"`
}

// Execute creates and persists the note, then broadcasts the "nota_clinica" event.
func (uc *AddNotaClinicaUseCase) Execute(
	pacienteID string,
	doctorID string,
	nombreDoctor string,
	nombrePaciente string,
	req AddNotaRequest,
) (*entities.PatientNote, error) {
	if req.Content == "" {
		return nil, errors.New("el contenido de la nota no puede estar vacío")
	}

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return nil, errors.New("paciente_id inválido")
	}
	docObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, errors.New("doctor_id inválido")
	}

	nota := &entities.PatientNote{
		ID:          primitive.NewObjectID(),
		PacienteID:  pacObjID,
		DoctorID:    docObjID,
		Content:     req.Content,
		CreatedDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	saved, err := uc.repo.AddNota(nota)
	if err != nil {
		return nil, errors.New("no se pudo guardar la nota: " + err.Error())
	}

	// ─── Domain event ────────────────────────────────────────────────────────
	// The full note is embedded in the event so the WebSocket consumer can
	// render it on connected clients without a secondary DB query.
	_ = uc.publisher.Publish(entities.DoctorEvent{
		Tipo:           "nota_clinica",
		PacienteID:     pacienteID,
		NombrePaciente: nombrePaciente,
		DoctorID:       doctorID,
		NombreDoctor:   nombreDoctor,
		Nota:           saved,
	})

	return saved, nil
}
