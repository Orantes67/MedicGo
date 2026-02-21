package routers

import (
	"ApiMedicGO/src/feature/login/infraestructure/dependencies_login"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupLoginRoutes registra las rutas del feature login.
func SetupLoginRoutes(router *gin.RouterGroup, db *mongo.Database) {
	controller := dependencies_login.SetupLoginDependencies(db)
	router.POST("/login", controller.Login)
}
