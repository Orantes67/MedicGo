package adapters

import (
	"context"
	"errors"
	"time"

	"ApiMedicGO/src/feature/enfermeros/domain/entities"
	doctorEntities "ApiMedicGO/src/feature/doctores/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoEnfermeroRepository satisfies EnfermeroRepository by querying the
// pacientes, notas_clinicas and users collections from the nurse's perspective.
type MongoEnfermeroRepository struct {
	collection *mongo.Collection
	notas      *mongo.Collection
	users      *mongo.Collection
}

func NewMongoEnfermeroRepository(db *mongo.Database) *MongoEnfermeroRepository {
	return &MongoEnfermeroRepository{
		collection: db.Collection("pacientes"),
		notas:      db.Collection("notas_clinicas"),
		users:      db.Collection("users"),
	}
}

// GetPacientesAsignados returns the read-model list of patients assigned to the nurse.
func (r *MongoEnfermeroRepository) GetPacientesAsignados(enfermeroID string) ([]*entities.PacienteResumen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(enfermeroID)
	if err != nil {
		return nil, errors.New("enfermero_id inválido")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"enfermero_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Raw documents decoded into a minimal struct that matches the pacientes collection.
	type rawPaciente struct {
		ID                  primitive.ObjectID  `bson:"_id"`
		Nombre              string              `bson:"nombre"`
		Apellido            string              `bson:"apellido"`
		Edad                int                 `bson:"edad"`
		AreaNombre          string              `bson:"area_nombre"`
		EstadoActual        string              `bson:"estado_actual"`
		NotaCondicion       string              `bson:"nota_condicion"`
		UltimaActualizacion string              `bson:"ultima_actualizacion"`
		DoctorID            *primitive.ObjectID `bson:"doctor_id,omitempty"`
	}

	var results []*entities.PacienteResumen
	for cursor.Next(ctx) {
		var raw rawPaciente
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		// Resolve doctor name from the users collection.
		nombreDoctor := ""
		if raw.DoctorID != nil {
			docCtx, docCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer docCancel()
			var docUser struct {
				Name string `bson:"name"`
			}
			if err2 := r.users.FindOne(docCtx, bson.M{"_id": raw.DoctorID}).Decode(&docUser); err2 == nil {
				nombreDoctor = docUser.Name
			}
		}

		results = append(results, &entities.PacienteResumen{
			ID:                  raw.ID.Hex(),
			Nombre:              raw.Nombre,
			Apellido:            raw.Apellido,
			Edad:                raw.Edad,
			AreaNombre:          raw.AreaNombre,
			EstadoActual:        raw.EstadoActual,
			NotaCondicion:       raw.NotaCondicion,
			NombreDoctor:        nombreDoctor,
			UltimaActualizacion: raw.UltimaActualizacion,
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateEstadoPaciente atomically updates a patient's state, note, and timestamp.
// The enfermeroID is included in the filter so a nurse can ONLY update patients
// assigned to her — no extra authorization check needed at the use-case layer.
func (r *MongoEnfermeroRepository) UpdateEstadoPaciente(
	pacienteID string,
	enfermeroID string,
	nuevoEstado string,
	nota string,
	timestamp string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return errors.New("paciente_id inválido")
	}
	enfObjID, err := primitive.ObjectIDFromHex(enfermeroID)
	if err != nil {
		return errors.New("enfermero_id inválido")
	}

	filter := bson.M{
		"_id":         pacObjID,
		"enfermero_id": enfObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"estado_actual":        nuevoEstado,
			"nota_condicion":       nota,
			"ultima_actualizacion": timestamp,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("paciente no encontrado o no asignado a esta enfermera")
	}
	return nil
}

// ─── GetNotasByPaciente ──────────────────────────────────────────────────────

// GetNotasByPaciente returns all clinical notes for a patient.
// enfermeroID is used to verify the patient is assigned to this nurse first.
func (r *MongoEnfermeroRepository) GetNotasByPaciente(pacienteID string, enfermeroID string) ([]*doctorEntities.PatientNote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return nil, errors.New("paciente_id inválido")
	}
	enfObjID, err := primitive.ObjectIDFromHex(enfermeroID)
	if err != nil {
		return nil, errors.New("enfermero_id inválido")
	}

	// Verify the patient is actually assigned to this nurse.
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": pacObjID, "enfermero_id": enfObjID})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("paciente no encontrado o no asignado a esta enfermera")
	}

	// Fetch all notes regardless of who wrote them (doctor or nurse).
	cursor, err := r.notas.Find(ctx, bson.M{"paciente_id": pacObjID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notas []*doctorEntities.PatientNote
	if err := cursor.All(ctx, &notas); err != nil {
		return nil, err
	}
	return notas, nil
}

// ─── AddNota ─────────────────────────────────────────────────────────────────

// AddNota persists a new clinical note written by the nurse.
// PatientNote.DoctorID holds the enfermero's ObjectID as the note author.
func (r *MongoEnfermeroRepository) AddNota(nota *doctorEntities.PatientNote) (*doctorEntities.PatientNote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.notas.InsertOne(ctx, nota)
	if err != nil {
		return nil, err
	}
	return nota, nil
}
