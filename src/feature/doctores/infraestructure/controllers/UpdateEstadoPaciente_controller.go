package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/doctores/application"

	"github.com/gin-gonic/gin"
)

type UpdateEstadoPacienteDoctorController struct {
	useCase *application.UpdateEstadoPacienteDoctorUseCase
}

func NewUpdateEstadoPacienteDoctorController(uc *application.UpdateEstadoPacienteDoctorUseCase) *UpdateEstadoPacienteDoctorController {
	return &UpdateEstadoPacienteDoctorController{useCase: uc}
}

// UpdateEstadoPaciente godoc
// @Summary      Actualizar estado clínico del paciente
// @Description  El doctor autenticado cambia el estado del paciente.
//               Dispara un DoctorEvent "estado_actualizado" listo para el WebSocket.
// @Tags         doctores
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                                   true  "ID del paciente"
// @Param        body  body  application.UpdateEstadoDoctorRequest    true  "Nuevo estado"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /api/v1/doctores/pacientes/{id}/estado [patch]
func (c *UpdateEstadoPacienteDoctorController) UpdateEstado(ctx *gin.Context) {
	pacienteID := ctx.Param("id")
	doctorID, _      := ctx.Get("user_id")
	licenseNumber, _ := ctx.Get("license_number")

	var req application.UpdateEstadoDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// X-Paciente-Nombre enriches the event payload; optional from the client.
	nombrePaciente := ctx.GetHeader("X-Paciente-Nombre")

	if err := c.useCase.Execute(
		pacienteID,
		doctorID.(string),
		licenseNumber.(string),
		nombrePaciente,
		req,
	); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "paciente no encontrado o no asignado a este doctor" {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estado actualizado exitosamente"})
}
