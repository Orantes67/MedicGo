package routers

import (
	"ApiMedicGO/src/core/middleware"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
	"ApiMedicGO/src/feature/pacientes/infraestructure/dependencies_paciente"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupPacientesRoutes registra las rutas del feature pacientes.
//
// Reglas de acceso:
//   - POST   /pacientes              → solo administrador (crear un nuevo paciente)
//   - GET    /pacientes              → todos (administrador ve todo; doctor/enfermero solo los suyos)
//   - GET    /pacientes/:id          → todos los autenticados
//   - PATCH  /pacientes/:id/asignar  → administrador (asigna doctor y/o enfermero)
//   - PUT    /pacientes/:id          → todos los autenticados (actualizar signos vitales)
//   - DELETE /pacientes/:id          → solo administrador
func SetupPacientesRoutes(router *gin.RouterGroup, db *mongo.Database) {
	ctrl := dependencies_paciente.SetupPacienteDependencies(db)

	pacientes := router.Group("/pacientes")
	pacientes.Use(middleware.AuthMiddleware())
	{
		// Todos los autenticados
		pacientes.GET("", ctrl.GetAll.GetAll)
		pacientes.GET("/:id", ctrl.GetByID.GetByID)
		pacientes.PUT("/:id", ctrl.Update.Update)

		// Solo administrador puede crear, asignar o eliminar
		adminOnly := pacientes.Group("")
		adminOnly.Use(middleware.RequireRole(
			loginEntities.RoleAdmin,
		))
		{
			adminOnly.POST("", ctrl.Create.Create)
			adminOnly.PATCH("/:id/asignar", ctrl.Assign.Assign)
			adminOnly.DELETE("/:id", ctrl.Delete.Delete)
		}
	}
}
