package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type CreatePacienteController struct {
	useCase *application.CreatePacienteUseCase
}

func NewCreatePacienteController(uc *application.CreatePacienteUseCase) *CreatePacienteController {
	return &CreatePacienteController{useCase: uc}
}

// Create godoc
// @Summary      Crear paciente
// @Description  Solo el administrador puede registrar un nuevo paciente
// @Tags         pacientes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  application.CreatePacienteRequest  true  "Datos del paciente"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/pacientes [post]
func (c *CreatePacienteController) Create(ctx *gin.Context) {
	var req application.CreatePacienteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paciente, err := c.useCase.Execute(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Paciente creado exitosamente",
		"paciente": paciente,
	})
}
