package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/register/application"

	"github.com/gin-gonic/gin"
)

// RegisterController maneja las peticiones HTTP del registro.
type RegisterController struct {
	useCase *application.RegisterUseCase
}

func NewRegisterController(useCase *application.RegisterUseCase) *RegisterController {
	return &RegisterController{useCase: useCase}
}

// Register godoc
// @Summary      Registrar usuario
// @Description  Crea un nuevo usuario en el sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  application.RegisterRequest  true  "Datos del usuario"
// @Success      201  {object}  application.RegisterResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /api/v1/register [post]
func (c *RegisterController) Register(ctx *gin.Context) {
	var req application.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := c.useCase.Execute(req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
