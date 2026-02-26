package adapters

import (
	"context"
	"errors"
	"time"

	adminEntities "ApiMedicGO/src/feature/admin/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAdminRepository implementa repositories.AdminRepository sobre MongoDB.
type MongoAdminRepository struct {
	pacientes *mongo.Collection
	users     *mongo.Collection
}

func NewMongoAdminRepository(db *mongo.Database) *MongoAdminRepository {
	return &MongoAdminRepository{
		pacientes: db.Collection("pacientes"),
		users:     db.Collection("users"),
	}
}

// ---------------------------------------------------------------------------
// GetMetrics – métricas del dashboard principal del administrador
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) GetMetrics() (*adminEntities.AdminMetrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- Conteo de pacientes agrupados por estado ---
	type estadoCount struct {
		Estado string `bson:"_id"`
		Total  int    `bson:"total"`
	}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$estado_actual"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}
	cursor, err := r.pacientes.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var estadosCounts []estadoCount
	if err := cursor.All(ctx, &estadosCounts); err != nil {
		return nil, err
	}

	metrics := &adminEntities.AdminMetrics{}
	for _, ec := range estadosCounts {
		metrics.TotalPacientes += ec.Total
		switch ec.Estado {
		case "critico":
			metrics.Criticos = ec.Total
		case "observacion":
			metrics.Observacion = ec.Total
		case "estable":
			metrics.Estables = ec.Total
		}
	}
	metrics.Alertas = metrics.Criticos // emergencias activas

	// --- Conteo de personal por rol ---
	countCtx, countCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer countCancel()

	totalEnfermeros, err := r.users.CountDocuments(countCtx, bson.M{"role": "enfermero"})
	if err != nil {
		return nil, err
	}
	totalDoctores, err := r.users.CountDocuments(countCtx, bson.M{"role": "doctor"})
	if err != nil {
		return nil, err
	}
	metrics.Enfermeros = int(totalEnfermeros)
	metrics.Doctores = int(totalDoctores)
	metrics.PersonalActivo = metrics.Enfermeros + metrics.Doctores

	// --- Actividad reciente: últimos 5 pacientes actualizados ---
	actCtx, actCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer actCancel()

	opts := options.Find().
		SetSort(bson.D{{Key: "ultima_actualizacion", Value: -1}}).
		SetLimit(5)

	actCursor, err := r.pacientes.Find(actCtx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer actCursor.Close(actCtx)

	type rawPaciente struct {
		Nombre          string `bson:"nombre"`
		Apellido        string `bson:"apellido"`
		AreaNombre      string `bson:"area_nombre"`
		NombreEnfermero string `bson:"nombre_enfermero"`
		EstadoActual    string `bson:"estado_actual"`
	}
	for actCursor.Next(actCtx) {
		var raw rawPaciente
		if err := actCursor.Decode(&raw); err != nil {
			continue
		}
		metrics.ActividadReciente = append(metrics.ActividadReciente, &adminEntities.ActividadItem{
			NombrePaciente:  raw.Nombre + " " + raw.Apellido,
			Area:            raw.AreaNombre,
			NombreEnfermero: raw.NombreEnfermero,
			Estado:          raw.EstadoActual,
		})
	}

	return metrics, nil
}

// ---------------------------------------------------------------------------
// GetUsuariosByRol – lista de usuarios por rol con su doctor asignado (enfermeros)
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) GetUsuariosByRol(rol string) ([]*adminEntities.UsuarioResumen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.users.Find(ctx, bson.M{"role": rol})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type rawUser struct {
		ID            primitive.ObjectID  `bson:"_id"`
		Name          string              `bson:"name"`
		Specialty     string              `bson:"specialty"`
		Role          string              `bson:"role"`
		DoctorID      *primitive.ObjectID `bson:"doctor_id,omitempty"`
	}

	var results []*adminEntities.UsuarioResumen
	for cursor.Next(ctx) {
		var raw rawUser
		if err := cursor.Decode(&raw); err != nil {
			continue
		}
		u := &adminEntities.UsuarioResumen{
			ID:           raw.ID.Hex(),
			Nombre:       raw.Name,
			Especialidad: raw.Specialty,
			Rol:          raw.Role,
		}
		if raw.DoctorID != nil {
			u.DoctorID = raw.DoctorID.Hex()
			// Resolver nombre del doctor
			docCtx, docCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer docCancel()
			var doc struct {
				Name string `bson:"name"`
			}
			if err2 := r.users.FindOne(docCtx, bson.M{"_id": raw.DoctorID}).Decode(&doc); err2 == nil {
				u.DoctorNombre = doc.Name
			}
		}
		results = append(results, u)
	}
	return results, nil
}

// ---------------------------------------------------------------------------
// FindByLicenseNumber – busca usuario por número de colegiado
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) FindByLicenseNumber(licencia string) (*adminEntities.UsuarioResumen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var raw struct {
		ID            primitive.ObjectID `bson:"_id"`
		Name          string             `bson:"name"`
		Specialty     string             `bson:"specialty"`
		Role          string             `bson:"role"`
	}
	err := r.users.FindOne(ctx, bson.M{"license_number": licencia}).Decode(&raw)
	if err != nil {
		return nil, err
	}
	return &adminEntities.UsuarioResumen{
		ID:           raw.ID.Hex(),
		Nombre:       raw.Name,
		Especialidad: raw.Specialty,
		Rol:          raw.Role,
	}, nil
}

// ---------------------------------------------------------------------------
// SaveUsuario – crea un nuevo usuario (enfermero o doctor)
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) SaveUsuario(nombre, email, licencia, hashedPassword, rol, especialidad string) (*adminEntities.UsuarioResumen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newID := primitive.NewObjectID()
	doc := bson.M{
		"_id":            newID,
		"name":           nombre,
		"email":          email,
		"license_number": licencia,
		"password":       hashedPassword,
		"role":           rol,
		"specialty":      especialidad,
	}
	if _, err := r.users.InsertOne(ctx, doc); err != nil {
		return nil, err
	}
	return &adminEntities.UsuarioResumen{
		ID:           newID.Hex(),
		Nombre:       nombre,
		Especialidad: especialidad,
		Rol:          rol,
	}, nil
}

// ---------------------------------------------------------------------------
// AsignarEnfermeroADoctor – actualiza el campo doctor_id del enfermero
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) AsignarEnfermeroADoctor(enfermeroID, doctorID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	enfObjID, err := primitive.ObjectIDFromHex(enfermeroID)
	if err != nil {
		return errors.New("enfermero_id inválido")
	}
	docObjID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return errors.New("doctor_id inválido")
	}

	result, err := r.users.UpdateOne(
		ctx,
		bson.M{"_id": enfObjID, "role": "enfermero"},
		bson.M{"$set": bson.M{"doctor_id": docObjID}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("enfermero no encontrado")
	}
	return nil
}

// ---------------------------------------------------------------------------
// GetDistribucionAreas – distribución de pacientes por área del hospital
// ---------------------------------------------------------------------------

func (r *MongoAdminRepository) GetDistribucionAreas() ([]*adminEntities.AreaDistribucion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Todas las áreas definidas (incluye las que pueden tener 0 pacientes)
	allAreas := []string{"urgencias", "hospitalizacion", "uci", "cirugia", "pediatria"}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$area_nombre"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "criticos", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$estado_actual", "critico"}}},
						1,
						0,
					}},
				}},
			}},
		}}},
	}
	cursor, err := r.pacientes.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type rawArea struct {
		Area     string `bson:"_id"`
		Total    int    `bson:"total"`
		Criticos int    `bson:"criticos"`
	}
	rawMap := make(map[string]rawArea)
	for cursor.Next(ctx) {
		var ra rawArea
		if err := cursor.Decode(&ra); err != nil {
			continue
		}
		rawMap[ra.Area] = ra
	}

	// Combinar con la lista estática para incluir áreas sin pacientes
	var result []*adminEntities.AreaDistribucion
	for _, area := range allAreas {
		dist := &adminEntities.AreaDistribucion{Area: area}
		if raw, ok := rawMap[area]; ok {
			dist.Total = raw.Total
			dist.Criticos = raw.Criticos
		}
		result = append(result, dist)
	}
	return result, nil
}
