package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/doctores/application"

	"github.com/gin-gonic/gin"
)

type GetMisPacientesDoctorController struct {
	useCase *application.GetMisPacientesDoctorUseCase
}

func NewGetMisPacientesDoctorController(uc *application.GetMisPacientesDoctorUseCase) *GetMisPacientesDoctorController {
	return &GetMisPacientesDoctorController{useCase: uc}
}

// GetMisPacientes godoc
// @Summary      Pacientes asignados al doctor
// @Description  Retorna stats, sección "Prioridad" (críticos) y lista completa de pacientes
//               asignados al doctor autenticado.
// @Tags         doctores
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  entities.MisPacientesDoctorResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/doctores/mis-pacientes [get]
func (c *GetMisPacientesDoctorController) GetMisPacientes(ctx *gin.Context) {
	doctorID, _ := ctx.Get("user_id")

	response, err := c.useCase.Execute(doctorID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
