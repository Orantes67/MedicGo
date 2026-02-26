package main

import (
	"log"
	"os"

	"ApiMedicGO/src/core"
	doctorRouters "ApiMedicGO/src/feature/doctores/infraestructure/routers"
	enfermeroRouters "ApiMedicGO/src/feature/enfermeros/infraestructure/routers"
	loginRouters "ApiMedicGO/src/feature/login/infraestructure/routers"
	pacientesRouters "ApiMedicGO/src/feature/pacientes/infraestructure/routers"
	registerRouters "ApiMedicGO/src/feature/register/infraestructure/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using system environment variables")
	}

	// Validar variables de entorno obligatorias
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("❌ JWT_SECRET no está configurado. Define la variable de entorno antes de iniciar.")
	}
	if os.Getenv("MONGO_CREDENTIALS") == "" {
		log.Fatal("❌ MONGO_CREDENTIALS no está configurado. Define la variable de entorno antes de iniciar.")
	}

	// Conectar a MongoDB
	core.ConnectMongoDB()

	// Configurar router de Gin
	router := gin.Default()

	// Grupo base de la API
	api := router.Group("/api/v1")
	{
		loginRouters.SetupLoginRoutes(api, core.DB)
		registerRouters.SetupRegisterRoutes(api, core.DB)
		pacientesRouters.SetupPacientesRoutes(api, core.DB)
		// publisher = nil → NoopPublisher until you inject the WebSocket hub.
		enfermeroRouters.SetupEnfermerosRoutes(api, core.DB, nil)
		doctorRouters.SetupDoctoresRoutes(api, core.DB, nil)
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(" Server running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
