package routers

import (
	"ApiMedicGO/src/core/events"
	"ApiMedicGO/src/core/middleware"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
	"ApiMedicGO/src/feature/enfermeros/infraestructure/dependencies_enfermeros"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupEnfermerosRoutes registers all routes for the enfermeros feature.
//
// Routes:
//   GET    /enfermeros/mis-pacientes              → enfermera autenticada (Home screen)
//   PATCH  /enfermeros/pacientes/:id/estado       → solo enfermera (Actualizar Estado + Nota Rápida)
//   GET    /enfermeros/pacientes/:id/notas        → solo enfermera (ver notas clínicas)
//   POST   /enfermeros/pacientes/:id/notas        → solo enfermera (agregar nota clínica)
//
// publisher: inject your WebSocket EventPublisher here.
//            Pass nil to use the built-in no-op publisher (no real-time events).
func SetupEnfermerosRoutes(router *gin.RouterGroup, db *mongo.Database, publisher events.EventPublisher) {
	ctrl := dependencies_enfermeros.SetupEnfermeroDependencies(db, publisher)

	enfermeros := router.Group("/enfermeros")
	enfermeros.Use(middleware.AuthMiddleware())
	{
		// Only nurses can access their own section.
		nurseOnly := enfermeros.Group("")
		nurseOnly.Use(middleware.RequireRole(loginEntities.RoleNurse))
		{
			nurseOnly.GET("/mis-pacientes", ctrl.GetMisPacientes.GetMisPacientes)
			nurseOnly.PATCH("/pacientes/:id/estado", ctrl.UpdateEstadoPaciente.UpdateEstado)
			nurseOnly.GET("/pacientes/:id/notas", ctrl.GetNotas.GetNotas)
			nurseOnly.POST("/pacientes/:id/notas", ctrl.AddNota.AddNota)
		}
	}
}
