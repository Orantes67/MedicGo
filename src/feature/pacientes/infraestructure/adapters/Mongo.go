package adapters

import (
	"context"
	"errors"
	"time"

	"ApiMedicGO/src/feature/pacientes/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPacienteRepository struct {
	collection *mongo.Collection
}

func NewMongoPacienteRepository(db *mongo.Database) *MongoPacienteRepository {
	return &MongoPacienteRepository{
		collection: db.Collection("pacientes"),
	}
}

func (r *MongoPacienteRepository) Save(paciente *entities.Paciente) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, paciente)
	return err
}

func (r *MongoPacienteRepository) FindAll() ([]*entities.Paciente, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pacientes []*entities.Paciente
	if err := cursor.All(ctx, &pacientes); err != nil {
		return nil, err
	}
	return pacientes, nil
}

func (r *MongoPacienteRepository) FindByID(id string) (*entities.Paciente, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("id inválido")
	}

	var paciente entities.Paciente
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&paciente); err != nil {
		return nil, err
	}
	return &paciente, nil
}

func (r *MongoPacienteRepository) Update(id string, paciente *entities.Paciente) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id inválido")
	}

	update := bson.M{
		"$set": bson.M{
			"nombre":               paciente.Nombre,
			"apellido":             paciente.Apellido,
			"edad":                 paciente.Edad,
			"area_nombre":          paciente.AreaNombre,
			"estado_actual":        paciente.EstadoActual,
			"nota_condicion":       paciente.NotaCondicion,
			"tipo_sangre":          paciente.TipoSangre,
			"sintomas":             paciente.Sintomas,
			"ultima_actualizacion": paciente.UltimaActualizacion,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *MongoPacienteRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("id inválido")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *MongoPacienteRepository) FindByDoctor(doctorID string) ([]*entities.Paciente, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, errors.New("doctor_id inválido")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"doctor_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pacientes []*entities.Paciente
	if err := cursor.All(ctx, &pacientes); err != nil {
		return nil, err
	}
	return pacientes, nil
}

func (r *MongoPacienteRepository) FindByEnfermero(enfermeroID string) ([]*entities.Paciente, error) {
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

	var pacientes []*entities.Paciente
	if err := cursor.All(ctx, &pacientes); err != nil {
		return nil, err
	}
	return pacientes, nil
}

// Assign actualiza el doctor_id y/o el enfermero_id de un paciente.
// Si el puntero es nil, el campo no se toca. Si apunta a "", se limpia.
func (r *MongoPacienteRepository) Assign(pacienteID string, doctorID *string, enfermeroID *string, nombreEnfermero string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(pacienteID)
	if err != nil {
		return errors.New("paciente_id inválido")
	}

	setFields := bson.M{}

	if doctorID != nil {
		if *doctorID == "" {
			setFields["doctor_id"] = nil
		} else {
			docObjID, err := primitive.ObjectIDFromHex(*doctorID)
			if err != nil {
				return errors.New("doctor_id inválido")
			}
			setFields["doctor_id"] = docObjID
		}
	}

	if enfermeroID != nil {
		if *enfermeroID == "" {
			setFields["enfermero_id"] = nil
			setFields["nombre_enfermero"] = ""
		} else {
			enfObjID, err := primitive.ObjectIDFromHex(*enfermeroID)
			if err != nil {
				return errors.New("enfermero_id inválido")
			}
			setFields["enfermero_id"] = enfObjID
			if nombreEnfermero != "" {
				setFields["nombre_enfermero"] = nombreEnfermero
			}
		}
	}

	if len(setFields) == 0 {
		return nil
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": setFields})
	return err
}
