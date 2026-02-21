package routers

import (
	"ApiMedicGO/src/feature/register/infraestructure/dependencies_register"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRegisterRoutes registra las rutas del feature register.
func SetupRegisterRoutes(router *gin.RouterGroup, db *mongo.Database) {
	controller := dependencies_register.SetupRegisterDependencies(db)
	router.POST("/register", controller.Register)
}
