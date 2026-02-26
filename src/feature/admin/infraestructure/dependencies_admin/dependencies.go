package dependencies_admin

import (
	"ApiMedicGO/src/feature/admin/application"
	"ApiMedicGO/src/feature/admin/infraestructure/adapters"
	"ApiMedicGO/src/feature/admin/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// AdminControllers agrupa todos los controladores del feature admin.
type AdminControllers struct {
	GetMetrics     *controllers.GetMetricsController
	GetUsuarios    *controllers.GetUsuariosController
	CreateUsuario  *controllers.CreateUsuarioController
	AsignarEnfermero *controllers.AsignarEnfermeroController
	GetAreas       *controllers.GetAreasController
}

// SetupAdminDependencies ensambla el grafo de dependencias del feature admin.
func SetupAdminDependencies(db *mongo.Database) *AdminControllers {
	repo := adapters.NewMongoAdminRepository(db)

	return &AdminControllers{
		GetMetrics:       controllers.NewGetMetricsController(application.NewGetMetricsUseCase(repo)),
		GetUsuarios:      controllers.NewGetUsuariosController(application.NewGetUsuariosUseCase(repo)),
		CreateUsuario:    controllers.NewCreateUsuarioController(application.NewCreateUsuarioUseCase(repo)),
		AsignarEnfermero: controllers.NewAsignarEnfermeroController(application.NewAsignarEnfermeroUseCase(repo)),
		GetAreas:         controllers.NewGetAreasController(application.NewGetAreasUseCase(repo)),
	}
}
