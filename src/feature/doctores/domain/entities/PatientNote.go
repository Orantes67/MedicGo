package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// PatientNote is a clinical note written by a doctor for a specific patient.
// It is stored in the "notas_clinicas" collection.
//
// Mirrors the mobile data class:
//
//	data class PatientNote(
//	    val id: Long,
//	    val patientId: Long,
//	    val doctorId: Long,
//	    val content: String,
//	    val createdDate: String
//	)
type PatientNote struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PacienteID   primitive.ObjectID `bson:"paciente_id"   json:"patientId"`
	DoctorID     primitive.ObjectID `bson:"doctor_id"     json:"doctorId"`
	Content      string             `bson:"content"       json:"content"`
	CreatedDate  string             `bson:"created_date"  json:"createdDate"`
}
