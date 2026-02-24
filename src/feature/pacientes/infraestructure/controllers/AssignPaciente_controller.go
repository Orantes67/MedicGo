package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type AssignPacienteController struct {
	useCase *application.AssignPacienteUseCase
}

func NewAssignPacienteController(uc *application.AssignPacienteUseCase) *AssignPacienteController {
	return &AssignPacienteController{useCase: uc}
}

// Assign godoc
// @Summary      Asignar paciente
// @Description  jefe_doctor asigna a un doctor | jefe_enfermera asigna a una enfermera
// @Tags         pacientes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                           true  "ID del paciente"
// @Param        body  body  application.AssignPacienteRequest true  "Asignación"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/pacientes/{id}/asignar [patch]
func (c *AssignPacienteController) Assign(ctx *gin.Context) {
	id := ctx.Param("id")
	role, _ := ctx.Get("role")

	var req application.AssignPacienteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.useCase.Execute(id, req, role.(string)); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Paciente asignado exitosamente"})
}
