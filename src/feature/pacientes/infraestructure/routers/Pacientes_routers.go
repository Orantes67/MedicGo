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
//   - POST   /pacientes              → todos los usuarios autenticados (crear un nuevo paciente)
//   - GET    /pacientes              → todos (jefes ven todo; doctor/enfermero solo los suyos)
//   - GET    /pacientes/:id          → todos los autenticados
//   - PATCH  /pacientes/:id/asignar  → jefe_doctor (asigna doctor) | jefe_enfermera (asigna enfermero)
//   - PUT    /pacientes/:id          → todos los autenticados (actualizar signos vitales)
//   - DELETE /pacientes/:id          → solo jefes
func SetupPacientesRoutes(router *gin.RouterGroup, db *mongo.Database) {
	ctrl := dependencies_paciente.SetupPacienteDependencies(db)

	pacientes := router.Group("/pacientes")
	pacientes.Use(middleware.AuthMiddleware())
	{
		// Todos los autenticados
		pacientes.POST("", ctrl.Create.Create)
		pacientes.GET("", ctrl.GetAll.GetAll)
		pacientes.GET("/:id", ctrl.GetByID.GetByID)
		pacientes.PUT("/:id", ctrl.Update.Update)

		// Solo jefes pueden asignar o eliminar
		jefeOnly := pacientes.Group("")
		jefeOnly.Use(middleware.RequireRole(
			loginEntities.RoleJefeDoctor,
			loginEntities.RoleJefeEnfermera,
		))
		{
			jefeOnly.PATCH("/:id/asignar", ctrl.Assign.Assign)
			jefeOnly.DELETE("/:id", ctrl.Delete.Delete)
		}
	}
}
