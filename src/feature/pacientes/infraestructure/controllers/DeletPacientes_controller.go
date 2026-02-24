package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type DeletePacienteController struct {
	useCase *application.DeletePacienteUseCase
}

func NewDeletePacienteController(uc *application.DeletePacienteUseCase) *DeletePacienteController {
	return &DeletePacienteController{useCase: uc}
}

func (c *DeletePacienteController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.useCase.Execute(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Paciente eliminado exitosamente"})
}
