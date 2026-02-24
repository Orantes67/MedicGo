package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims refleja el payload guardado en el JWT del login.
type Claims struct {
	UserID        string `json:"user_id"`
	LicenseNumber string `json:"license_number"`
	Role          string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware valida el token JWT y guarda los claims en el contexto.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("license_number", claims.LicenseNumber)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

// RequireRole solo permite el acceso a los roles indicados.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
			return
		}

		for _, r := range roles {
			if role == r {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no tienes permisos para realizar esta acción",
		})
	}
}
