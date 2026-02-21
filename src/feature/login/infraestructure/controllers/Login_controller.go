package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/login/application"

	"github.com/gin-gonic/gin"
)

// LoginController maneja las peticiones HTTP del login.
type LoginController struct {
	useCase *application.LoginUseCase
}

func NewLoginController(useCase *application.LoginUseCase) *LoginController {
	return &LoginController{useCase: useCase}
}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica un usuario y retorna un JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  application.LoginRequest  true  "Credenciales"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/login [post]
func (c *LoginController) Login(ctx *gin.Context) {
	var req application.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := c.useCase.Execute(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"token":   resp.Token,
		"user":    resp.User,
	})
}
