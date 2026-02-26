package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/admin/application"

	"github.com/gin-gonic/gin"
)

// CreateUsuarioController gestiona la creación de enfermeros y doctores por el admin.
type CreateUsuarioController struct {
	useCase *application.CreateUsuarioUseCase
}

func NewCreateUsuarioController(uc *application.CreateUsuarioUseCase) *CreateUsuarioController {
	return &CreateUsuarioController{useCase: uc}
}

// Create godoc
// @Summary      Crear usuario (enfermero o doctor)
// @Description  Solo el administrador puede registrar nuevos enfermeros o doctores.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  application.CreateUsuarioRequest  true  "Datos del usuario"
// @Success      201  {object}  entities.UsuarioResumen
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /api/v1/admin/usuarios [post]
func (c *CreateUsuarioController) Create(ctx *gin.Context) {
	var req application.CreateUsuarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario, err := c.useCase.Execute(req)
	if err != nil {
		if err.Error() == "el número de colegiado ya está registrado" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, usuario)
}
