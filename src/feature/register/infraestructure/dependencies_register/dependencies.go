package dependencies_register

import (
	"ApiMedicGO/src/feature/register/application"
	"ApiMedicGO/src/feature/register/infraestructure/adapters"
	"ApiMedicGO/src/feature/register/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRegisterDependencies inyecta las dependencias del feature register
// y retorna el controlador listo para usar.
func SetupRegisterDependencies(db *mongo.Database) *controllers.RegisterController {
	repo := adapters.NewMongoUserRepository(db)
	useCase := application.NewRegisterUseCase(repo)
	controller := controllers.NewRegisterController(useCase)
	return controller
}
