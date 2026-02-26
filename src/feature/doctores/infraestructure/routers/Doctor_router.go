package routers

import (
	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/core/middleware"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
	"ApiMedicGO/src/feature/doctores/infraestructure/dependencies_doctores"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupDoctoresRoutes registers all routes for the doctores feature.
//
// Routes:
//   GET    /doctores/mis-pacientes              → doctor autenticado (Home screen)
//   GET    /doctores/pacientes/:id              → doctor autenticado (Cards detail screen)
//   PATCH  /doctores/pacientes/:id/estado       → solo doctor (Actualizar Estado)
//   POST   /doctores/pacientes/:id/notas        → solo doctor (Nota Clínica)
//
// publisher: inject your WebSocket EventPublisher here.
//            Pass nil to run with the built-in no-op publisher.
func SetupDoctoresRoutes(router *gin.RouterGroup, db *mongo.Database, publisher events.EventPublisher) {
	ctrl := dependencies_doctores.SetupDoctorDependencies(db, publisher)

	doctores := router.Group("/doctores")
	doctores.Use(middleware.AuthMiddleware())
	{
		// Doctor-only endpoints
		doctorOnly := doctores.Group("")
		doctorOnly.Use(middleware.RequireRole(loginEntities.RoleDoctor))
		{
			doctorOnly.GET("/mis-pacientes", ctrl.GetMisPacientes.GetMisPacientes)
			doctorOnly.GET("/pacientes/:id", ctrl.GetPacienteDetalle.GetPacienteDetalle)
			doctorOnly.PATCH("/pacientes/:id/estado", ctrl.UpdateEstadoPaciente.UpdateEstado)
			doctorOnly.POST("/pacientes/:id/notas", ctrl.AddNotaClinica.AddNota)
		}
	}
}
