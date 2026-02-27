package main

import (
	"log"
	"net/http"
	"os"

	"ApiMedicGO/src/core"
	coreWs "ApiMedicGO/src/core/websocket"
	adminRouters "ApiMedicGO/src/feature/admin/infraestructure/routers"
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
		log.Println("No .env file found, using system environment variables")
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

	// Hub WebSocket (implementa EventPublisher)
	hub := coreWs.NewHub()

	// Configurar router de Gin
	router := gin.Default()

	// ── WebSocket ─────────────────────────────────────────────────────────────
	// wss://<host>/ws?token=<JWT>
	router.GET("/ws", coreWs.WsHandler(hub))

	// ── Health check ──────────────────────────────────────────────────────────
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── Auth (login / register) ───────────────────────────────────────────────
	auth := router.Group("/api/auth")
	{
		loginRouters.SetupLoginRoutes(auth, core.DB)
		registerRouters.SetupRegisterRoutes(auth, core.DB)
	}

	// ── Resto de la API ───────────────────────────────────────────────────────
	api := router.Group("/api/v1")
	{
		pacientesRouters.SetupPacientesRoutes(api, core.DB)
		enfermeroRouters.SetupEnfermerosRoutes(api, core.DB, hub)
		doctorRouters.SetupDoctoresRoutes(api, core.DB, hub)
		adminRouters.SetupAdminRoutes(api, core.DB)
	}

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
