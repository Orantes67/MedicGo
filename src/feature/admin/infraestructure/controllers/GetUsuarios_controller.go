package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/admin/application"

	"github.com/gin-gonic/gin"
)

// GetUsuariosController maneja la solicitud de listado de usuarios (tab Usuarios).
type GetUsuariosController struct {
	useCase *application.GetUsuariosUseCase
}

func NewGetUsuariosController(uc *application.GetUsuariosUseCase) *GetUsuariosController {
	return &GetUsuariosController{useCase: uc}
}

// GetUsuarios godoc
// @Summary      Listar usuarios del sistema
// @Description  Devuelve enfermeros y doctores separados en dos listas.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  application.GetUsuariosResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/admin/usuarios [get]
func (c *GetUsuariosController) GetUsuarios(ctx *gin.Context) {
	resp, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
