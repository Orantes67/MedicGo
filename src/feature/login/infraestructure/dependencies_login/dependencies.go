package dependencies_login

import (
	"ApiMedicGO/src/feature/login/application"
	"ApiMedicGO/src/feature/login/infraestructure/adapters"
	"ApiMedicGO/src/feature/login/infraestructure/controllers"

	"go.mongodb.org/mongo-driver/mongo"
)

// SetupLoginDependencies inyecta las dependencias del feature login
// y retorna el controlador listo para usar.
func SetupLoginDependencies(db *mongo.Database) *controllers.LoginController {
	repo := adapters.NewMongoUserRepository(db)
	useCase := application.NewLoginUseCase(repo)
	controller := controllers.NewLoginController(useCase)
	return controller
}
