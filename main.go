package main

import (
	"log"
	"os"

	"ApiMedicGO/src/core"
	loginRouters "ApiMedicGO/src/feature/login/infraestructure/routers"
	registerRouters "ApiMedicGO/src/feature/register/infraestructure/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using system environment variables")
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
