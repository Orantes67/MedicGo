package adapters

import (
	"context"
	"errors"
	"time"

	"ApiMedicGO/src/feature/doctores/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDoctorRepository satisfies DoctorRepository using MongoDB.
// It queries the "pacientes", "notas_clinicas" and "users" collections.
type MongoDoctorRepository struct {
	pacientes *mongo.Collection
	notas     *mongo.Collection
	users     *mongo.Collection
}

func NewMongoDoctorRepository(db *mongo.Database) *MongoDoctorRepository {
	return &MongoDoctorRepository{
		pacientes: db.Collection("pacientes"),
		notas:     db.Collection("notas_clinicas"),
		users:     db.Collection("users"),
	}
}

// ─── GetPacientesAsignados ───────────────────────────────────────────────────

func (r *MongoDoctorRepository) GetPacientesAsignados(doctorID string) ([]*entities.PacienteResumenDoctor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, errors.New("doctor_id inválido")
	}

	cursor, err := r.pacientes.Find(ctx, bson.M{"doctor_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type rawPaciente struct {
		ID                  primitive.ObjectID  `bson:"_id"`
		Nombre              string              `bson:"nombre"`
		Apellido            string              `bson:"apellido"`
		Edad                int                 `bson:"edad"`
		AreaNombre          string              `bson:"area_nombre"`
		EstadoActual        string              `bson:"estado_actual"`
		NotaCondicion       string              `bson:"nota_condicion"`
		NombreEnfermero     string              `bson:"nombre_enfermero"`
		UltimaActualizacion string              `bson:"ultima_actualizacion"`
		FechaRegistro       string              `bson:"fecha_registro"`
		EnfermeroID         *primitive.ObjectID `bson:"enfermero_id,omitempty"`
	}

	var results []*entities.PacienteResumenDoctor
	for cursor.Next(ctx) {
		var raw rawPaciente
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		// Resolve nurse name: prefer the denormalized field, fallback to users lookup.
		nombreEnfermero := raw.NombreEnfermero
		if nombreEnfermero == "" && raw.EnfermeroID != nil {
			enfCtx, enfCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer enfCancel()
			var enfUser struct {
				Name string `bson:"name"`
			}
			if err2 := r.users.FindOne(enfCtx, bson.M{"_id": raw.EnfermeroID}).Decode(&enfUser); err2 == nil {
				nombreEnfermero = enfUser.Name
			}
		}

		results = append(results, &entities.PacienteResumenDoctor{
			ID:                  raw.ID.Hex(),
			Nombre:              raw.Nombre,
			Apellido:            raw.Apellido,
			Edad:                raw.Edad,
			AreaNombre:          raw.AreaNombre,
			EstadoActual:        raw.EstadoActual,
			NotaCondicion:       raw.NotaCondicion,
			NombreEnfermero:     nombreEnfermero,
			UltimaActualizacion: raw.UltimaActualizacion,
			FechaRegistro:       raw.FechaRegistro,
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// ─── GetPacienteDetalle ──────────────────────────────────────────────────────

func (r *MongoDoctorRepository) GetPacienteDetalle(pacienteID string, doctorID string) (*entities.PacienteDetalleDoctor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return nil, errors.New("paciente_id inválido")
	}
	docObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, errors.New("doctor_id inválido")
	}

	type rawPaciente struct {
		ID                  primitive.ObjectID  `bson:"_id"`
		Nombre              string              `bson:"nombre"`
		Apellido            string              `bson:"apellido"`
		Edad                int                 `bson:"edad"`
		AreaNombre          string              `bson:"area_nombre"`
		EstadoActual        string              `bson:"estado_actual"`
		NombreEnfermero     string              `bson:"nombre_enfermero"`
		UltimaActualizacion string              `bson:"ultima_actualizacion"`
		EnfermeroID         *primitive.ObjectID `bson:"enfermero_id,omitempty"`
	}

	var raw rawPaciente
	err = r.pacientes.FindOne(ctx, bson.M{
		"_id":       pacObjID,
		"doctor_id": docObjID,
	}).Decode(&raw)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("paciente no encontrado o no asignado a este doctor")
		}
		return nil, err
	}

	// Resolve nurse name: prefer the denormalized field, fallback to users lookup.
	nombreEnfermero := raw.NombreEnfermero
	if nombreEnfermero == "" && raw.EnfermeroID != nil {
		enfCtx, enfCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer enfCancel()
		var enfUser struct {
			Name string `bson:"name"`
		}
		if err2 := r.users.FindOne(enfCtx, bson.M{"_id": raw.EnfermeroID}).Decode(&enfUser); err2 == nil {
			nombreEnfermero = enfUser.Name
		}
	}

	// Fetch clinical notes written by this doctor for this patient.
	notas, err := r.GetNotasByPaciente(pacienteID, doctorID)
	if err != nil {
		return nil, err
	}

	return &entities.PacienteDetalleDoctor{
		ID:                  raw.ID.Hex(),
		Nombre:              raw.Nombre,
		Apellido:            raw.Apellido,
		Edad:                raw.Edad,
		AreaNombre:          raw.AreaNombre,
		EstadoActual:        raw.EstadoActual,
		NombreEnfermero:     nombreEnfermero,
		UltimaActualizacion: raw.UltimaActualizacion,
		NotasClinicas:       notas,
	}, nil
}

// ─── UpdateEstadoPaciente ────────────────────────────────────────────────────

func (r *MongoDoctorRepository) UpdateEstadoPaciente(
	pacienteID string,
	doctorID string,
	nuevoEstado string,
	timestamp string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return errors.New("paciente_id inválido")
	}
	docObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return errors.New("doctor_id inválido")
	}

	result, err := r.pacientes.UpdateOne(ctx,
		bson.M{"_id": pacObjID, "doctor_id": docObjID},
		bson.M{"$set": bson.M{
			"estado_actual":        nuevoEstado,
			"ultima_actualizacion": timestamp,
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("paciente no encontrado o no asignado a este doctor")
	}
	return nil
}

// ─── AddNota ────────────────────────────────────────────────────────────────

func (r *MongoDoctorRepository) AddNota(nota *entities.PatientNote) (*entities.PatientNote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.notas.InsertOne(ctx, nota)
	if err != nil {
		return nil, err
	}
	return nota, nil
}

// ─── GetNotasByPaciente ──────────────────────────────────────────────────────

func (r *MongoDoctorRepository) GetNotasByPaciente(pacienteID string, doctorID string) ([]*entities.PatientNote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pacObjID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return nil, errors.New("paciente_id inválido")
	}
	docObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, errors.New("doctor_id inválido")
	}

	cursor, err := r.notas.Find(ctx, bson.M{
		"paciente_id": pacObjID,
		"doctor_id":   docObjID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notas []*entities.PatientNote
	if err := cursor.All(ctx, &notas); err != nil {
		return nil, err
	}
	return notas, nil
}
