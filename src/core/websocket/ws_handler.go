package websocket

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permitir cualquier origen. Ajusta según tus necesidades de seguridad.
		return true
	},
}

// Claims refleja el payload del JWT (mismo que AuthMiddleware).
type wsClaims struct {
	UserID        string `json:"user_id"`
	LicenseNumber string `json:"license_number"`
	Role          string `json:"role"`
	jwt.RegisteredClaims
}

// WsHandler devuelve un gin.HandlerFunc que:
//  1. Valida el token JWT enviado como query param ?token=<JWT>
//  2. Hace upgrade a WebSocket
//  3. Registra la conexión en el Hub
func WsHandler(hub *Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Query("token")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		claims := &wsClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			// upgrader ya escribió el error HTTP, solo registramos el log.
			return
		}

		hub.Register(conn)
	}
}
