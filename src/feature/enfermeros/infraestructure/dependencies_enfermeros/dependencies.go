package dependencies_enfermeros

import (
	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/feature/enfermeros/application"
	"ApiMedicGO/src/feature/enfermeros/infraestructure/adapters"
	"ApiMedicGO/src/feature/enfermeros/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// EnfermeroControllers groups all controllers for the enfermeros feature.
type EnfermeroControllers struct {
	GetMisPacientes      *controllers.GetMisPacientesController
	UpdateEstadoPaciente *controllers.UpdateEstadoPacienteController
	GetNotas             *controllers.GetNotasPacienteController
	AddNota              *controllers.AddNotaClinicaEnfermeroController
}

// SetupEnfermeroDependencies wires the full dependency graph for the enfermeros feature.
//
// publisher: pass your WebSocket hub once implemented, or nil to use the no-op default.
// Keeping the publisher optional means the API starts and works without a WebSocket
// server; real-time events will simply be discarded until a real publisher is injected.
func SetupEnfermeroDependencies(db *mongo.Database, publisher events.EventPublisher) *EnfermeroControllers {
	if publisher == nil {
		publisher = &events.NoopPublisher{}
	}

	repo := adapters.NewMongoEnfermeroRepository(db)

	return &EnfermeroControllers{
		GetMisPacientes:      controllers.NewGetMisPacientesController(application.NewGetMisPacientesUseCase(repo)),
		UpdateEstadoPaciente: controllers.NewUpdateEstadoPacienteController(application.NewUpdateEstadoPacienteUseCase(repo, publisher)),
		GetNotas:             controllers.NewGetNotasPacienteController(application.NewGetNotasPacienteUseCase(repo)),
		AddNota:              controllers.NewAddNotaClinicaEnfermeroController(application.NewAddNotaClinicaEnfermeroUseCase(repo, publisher)),
	}
}
