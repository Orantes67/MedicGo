package routers

import (
	"ApiMedicGO/src/core/middleware"
	loginEntities "ApiMedicGO/src/feature/login/domain/entities"
	"ApiMedicGO/src/feature/admin/infraestructure/dependencies_admin"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAdminRoutes registra todas las rutas del panel administrativo.
//
// Todas las rutas requieren autenticación JWT y rol "administrador".
//
// Rutas:
//   GET   /admin/metricas                        → métricas del dashboard (tab Métricas)
//   GET   /admin/usuarios                        → lista de enfermeros y doctores (tab Usuarios)
//   POST  /admin/usuarios                        → crear enfermero o doctor
//   PATCH /admin/usuarios/asignar-enfermero      → asignar enfermero a un doctor
//   GET   /admin/areas                           → distribución de pacientes por área (tab Áreas)
func SetupAdminRoutes(router *gin.RouterGroup, db *mongo.Database) {
	ctrl := dependencies_admin.SetupAdminDependencies(db)

	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRole(loginEntities.RoleAdmin))
	{
		admin.GET("/metricas", ctrl.GetMetrics.GetMetrics)
		admin.GET("/usuarios", ctrl.GetUsuarios.GetUsuarios)
		admin.POST("/usuarios", ctrl.CreateUsuario.Create)
		admin.PATCH("/usuarios/asignar-enfermero", ctrl.AsignarEnfermero.Asignar)
		admin.GET("/areas", ctrl.GetAreas.GetAreas)
	}
}
