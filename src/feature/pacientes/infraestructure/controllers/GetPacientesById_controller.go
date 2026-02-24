package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type GetPacienteByIDController struct {
	useCase *application.GetPacienteByIDUseCase
}

func NewGetPacienteByIDController(uc *application.GetPacienteByIDUseCase) *GetPacienteByIDController {
	return &GetPacienteByIDController{useCase: uc}
}

func (c *GetPacienteByIDController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	paciente, err := c.useCase.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "paciente no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"paciente": paciente})
}
