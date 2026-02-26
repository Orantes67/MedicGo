package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Paciente struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre string             `bson:"nombre"        json:"nombre"`
	// Apellido del paciente
	Apellido string `bson:"apellido" json:"apellido"`
	Edad     int    `bson:"edad"     json:"edad"`

	// Área hospitalaria
	AreaID     *primitive.ObjectID `bson:"area_id,omitempty"  json:"area_id,omitempty"`
	AreaNombre string              `bson:"area_nombre"        json:"area_nombre"`

	// Estado clínico
	EstadoActual  string `bson:"estado_actual"  json:"estado_actual"`
	NotaCondicion string `bson:"nota_condicion" json:"nota_condicion"`

	// Datos médicos básicos
	TipoSangre    string `bson:"tipo_sangre"    json:"tipo_sangre"`
	Sintomas      string `bson:"sintomas"       json:"sintomas"`
	FechaRegistro string `bson:"fecha_registro" json:"fecha_registro"`

	// Asignaciones
	DoctorID    *primitive.ObjectID `bson:"doctor_id,omitempty"    json:"doctor_id,omitempty"`
	EnfermeroID *primitive.ObjectID `bson:"enfermero_id,omitempty" json:"enfermero_id,omitempty"`
	// Nombre denormalizado del enfermero para consultas rápidas
	NombreEnfermero string `bson:"nombre_enfermero,omitempty" json:"nombre_enfermero,omitempty"`

	// Control de tiempo
	UltimaActualizacion string `bson:"ultima_actualizacion" json:"ultima_actualizacion"`
}
