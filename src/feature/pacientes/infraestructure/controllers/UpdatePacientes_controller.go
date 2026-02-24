package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type UpdatePacienteController struct {
	useCase *application.UpdatePacienteUseCase
}

func NewUpdatePacienteController(uc *application.UpdatePacienteUseCase) *UpdatePacienteController {
	return &UpdatePacienteController{useCase: uc}
}

func (c *UpdatePacienteController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req application.UpdatePacienteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paciente, err := c.useCase.Execute(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Paciente actualizado exitosamente",
		"paciente": paciente,
	})
}
