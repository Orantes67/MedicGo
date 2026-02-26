package application

import (
	"errors"
	"time"

	"ApiMedicGO/src/core/events"
	doctorEntities "ApiMedicGO/src/feature/doctores/domain/entities"
	"ApiMedicGO/src/feature/enfermeros/domain/entities"
	"ApiMedicGO/src/feature/enfermeros/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddNotaClinicaEnfermeroUseCase lets a nurse write a clinical note for one of her
// assigned patients.  After persisting it fires a PacienteEstadoEvent with
// tipo = "nota_clinica" so the WebSocket hub can broadcast in real time.
//
// The PatientNote.DoctorID field stores the nurse's ObjectID as the note author
// (same collection, same structure — the "doctorId" field holds any medical-staff author).
type AddNotaClinicaEnfermeroUseCase struct {
	repo      repositories.EnfermeroRepository
	publisher events.EventPublisher
}

func NewAddNotaClinicaEnfermeroUseCase(
	repo repositories.EnfermeroRepository,
	publisher events.EventPublisher,
) *AddNotaClinicaEnfermeroUseCase {
	return &AddNotaClinicaEnfermeroUseCase{repo: repo, publisher: publisher}
}

// AddNotaEnfermeroRequest is the DTO sent by the nurse.
type AddNotaEnfermeroRequest struct {
	Content string `json:"content" binding:"required"`
}

// Execute creates, persists and broadcasts the clinical note.
func (uc *AddNotaClinicaEnfermeroUseCase) Execute(
	pacienteID string,
	enfermeroID string,
	nombreEnfermero string,
	nombrePaciente string,
	req AddNotaEnfermeroRequest,
) (*doctorEntities.PatientNote, error) {
	if req.Content == "" {
		return nil, errors.New("el contenido de la nota no puede estar vacío")
	}

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return nil, errors.New("paciente_id inválido")
	}
	enfObjID, err := primitive.ObjectIDFromHex(enfermeroID)
	if err != nil {
		return nil, errors.New("enfermero_id inválido")
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	nota := &doctorEntities.PatientNote{
		ID:          primitive.NewObjectID(),
		PacienteID:  pacObjID,
		DoctorID:    enfObjID, // nurse is the author; field reused per shared data class
		Content:     req.Content,
		CreatedDate: now,
	}

	saved, err := uc.repo.AddNota(nota)
	if err != nil {
		return nil, errors.New("no se pudo guardar la nota: " + err.Error())
	}

	// ─── Domain event ────────────────────────────────────────────────────────
	_ = uc.publisher.Publish(entities.PacienteEstadoEvent{
		Tipo:            "nota_clinica",
		PacienteID:      pacienteID,
		NombrePaciente:  nombrePaciente,
		EnfermeroID:     enfermeroID,
		NombreEnfermero: nombreEnfermero,
		NoteID:          saved.ID.Hex(),
		NoteContent:     saved.Content,
		NoteDate:        saved.CreatedDate,
	})

	return saved, nil
}
