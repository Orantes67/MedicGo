package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Paciente struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"       json:"id"`
	Nombre          string              `bson:"nombre"              json:"nombre"`
	Edad            int                 `bson:"edad"                json:"edad"`
	Estatura        float64             `bson:"estatura"            json:"estatura"`          // cm
	Peso            float64             `bson:"peso"                json:"peso"`              // kg
	RitmoCardiaco   int                 `bson:"ritmo_cardiaco"      json:"ritmo_cardiaco"`    // BPM
	Saturacion      float64             `bson:"saturacion"          json:"saturacion"`        // %
	PresionArterial string              `bson:"presion_arterial"    json:"presion_arterial"`  // ej. "120/80"
	SpO2            float64             `bson:"spo2"                json:"spo2"`              // %
	// Asignaciones
	DoctorID    *primitive.ObjectID `bson:"doctor_id,omitempty"    json:"doctor_id,omitempty"`
	EnfermeroID *primitive.ObjectID `bson:"enfermero_id,omitempty" json:"enfermero_id,omitempty"`
}
