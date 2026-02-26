package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/doctores/application"

	"github.com/gin-gonic/gin"
)

type GetPacienteDetalleController struct {
	useCase *application.GetPacienteDetalleUseCase
}

func NewGetPacienteDetalleController(uc *application.GetPacienteDetalleUseCase) *GetPacienteDetalleController {
	return &GetPacienteDetalleController{useCase: uc}
}

// GetPacienteDetalle godoc
// @Summary      Detalle del paciente (vista Cards)
// @Description  Retorna el modelo completo del paciente incluyendo estado, área,
//               enfermera asignada, última actualización y notas clínicas del doctor.
// @Tags         doctores
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "ID del paciente"
// @Success      200  {object}  entities.PacienteDetalleDoctor
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/doctores/pacientes/{id} [get]
func (c *GetPacienteDetalleController) GetPacienteDetalle(ctx *gin.Context) {
	pacienteID := ctx.Param("id")
	doctorID, _ := ctx.Get("user_id")

	detalle, err := c.useCase.Execute(pacienteID, doctorID.(string))
	if err != nil {
		status := http.StatusInternalServerError
		msg := err.Error()
		if msg == "paciente no encontrado o no asignado a este doctor" {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, detalle)
}
