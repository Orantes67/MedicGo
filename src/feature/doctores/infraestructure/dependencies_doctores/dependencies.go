package dependencies_doctores

import (
	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/feature/doctores/application"
	"ApiMedicGO/src/feature/doctores/infraestructure/adapters"
	"ApiMedicGO/src/feature/doctores/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// DoctorControllers groups all controllers for the doctores feature.
type DoctorControllers struct {
	GetMisPacientes      *controllers.GetMisPacientesDoctorController
	GetPacienteDetalle   *controllers.GetPacienteDetalleController
	UpdateEstadoPaciente *controllers.UpdateEstadoPacienteDoctorController
	AddNotaClinica       *controllers.AddNotaClinicaController
}

// SetupDoctorDependencies wires the full dependency graph for the doctores feature.
//
// publisher: pass your WebSocket hub once it is implemented, or nil to use NoopPublisher.
func SetupDoctorDependencies(db *mongo.Database, publisher events.EventPublisher) *DoctorControllers {
	if publisher == nil {
		publisher = &events.NoopPublisher{}
	}

	repo := adapters.NewMongoDoctorRepository(db)

	return &DoctorControllers{
		GetMisPacientes:      controllers.NewGetMisPacientesDoctorController(application.NewGetMisPacientesDoctorUseCase(repo)),
		GetPacienteDetalle:   controllers.NewGetPacienteDetalleController(application.NewGetPacienteDetalleUseCase(repo)),
		UpdateEstadoPaciente: controllers.NewUpdateEstadoPacienteDoctorController(application.NewUpdateEstadoPacienteDoctorUseCase(repo, publisher)),
		AddNotaClinica:       controllers.NewAddNotaClinicaController(application.NewAddNotaClinicaUseCase(repo, publisher)),
	}
}
