package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/admin/application"

	"github.com/gin-gonic/gin"
)

// AsignarEnfermeroController gestiona la asignación de un enfermero a un doctor.
type AsignarEnfermeroController struct {
	useCase *application.AsignarEnfermeroUseCase
}

func NewAsignarEnfermeroController(uc *application.AsignarEnfermeroUseCase) *AsignarEnfermeroController {
	return &AsignarEnfermeroController{useCase: uc}
}

// Asignar godoc
// @Summary      Asignar enfermero a un doctor
// @Description  El administrador vincula un enfermero a un doctor específico.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  application.AsignarEnfermeroRequest  true  "IDs del enfermero y del doctor"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/admin/usuarios/asignar-enfermero [patch]
func (c *AsignarEnfermeroController) Asignar(ctx *gin.Context) {
	var req application.AsignarEnfermeroRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.useCase.Execute(req); err != nil {
		if err.Error() == "enfermero no encontrado" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "enfermero asignado correctamente"})
}
