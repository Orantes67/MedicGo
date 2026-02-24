package dependencies_paciente

import (
	"ApiMedicGO/src/feature/pacientes/application"
	"ApiMedicGO/src/feature/pacientes/infraestructure/adapters"
	"ApiMedicGO/src/feature/pacientes/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// PacienteControllers agrupa todos los controladores del feature pacientes.
type PacienteControllers struct {
	Create  *controllers.CreatePacienteController
	GetAll  *controllers.GetPacientesController
	GetByID *controllers.GetPacienteByIDController
	Update  *controllers.UpdatePacienteController
	Delete  *controllers.DeletePacienteController
	Assign  *controllers.AssignPacienteController
}

// SetupPacienteDependencies inyecta todas las dependencias del feature pacientes.
func SetupPacienteDependencies(db *mongo.Database) *PacienteControllers {
	repo := adapters.NewMongoPacienteRepository(db)

	return &PacienteControllers{
		Create:  controllers.NewCreatePacienteController(application.NewCreatePacienteUseCase(repo)),
		GetAll:  controllers.NewGetPacientesController(application.NewGetPacientesUseCase(repo)),
		GetByID: controllers.NewGetPacienteByIDController(application.NewGetPacienteByIDUseCase(repo)),
		Update:  controllers.NewUpdatePacienteController(application.NewUpdatePacienteUseCase(repo)),
		Delete:  controllers.NewDeletePacienteController(application.NewDeletePacienteUseCase(repo)),
		Assign:  controllers.NewAssignPacienteController(application.NewAssignPacienteUseCase(repo)),
	}
}
