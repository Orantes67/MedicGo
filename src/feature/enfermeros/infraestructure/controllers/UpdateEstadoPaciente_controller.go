package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/enfermeros/application"

	"github.com/gin-gonic/gin"
)

type UpdateEstadoPacienteController struct {
	useCase *application.UpdateEstadoPacienteUseCase
}

func NewUpdateEstadoPacienteController(uc *application.UpdateEstadoPacienteUseCase) *UpdateEstadoPacienteController {
	return &UpdateEstadoPacienteController{useCase: uc}
}

// UpdateEstadoPaciente godoc
// @Summary      Actualizar estado del paciente
// @Description  La enfermera autenticada cambia el estado clínico y/o deja una nota rápida.
//               Internamente dispara un PacienteEstadoEvent listo para conectar al WebSocket.
// @Tags         enfermeros
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                               true  "ID del paciente"
// @Param        body  body  application.UpdateEstadoRequest      true  "Nuevo estado y nota"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/enfermeros/pacientes/{id}/estado [patch]
func (c *UpdateEstadoPacienteController) UpdateEstado(ctx *gin.Context) {
	pacienteID := ctx.Param("id")

	// Claims injected by AuthMiddleware
	enfermeroID, _    := ctx.Get("user_id")
	licenseNumber, _  := ctx.Get("license_number")

	// Use license_number as the display name that goes into the event.
	// When richer nurse data is available (e.g. full name stored in DB) this
	// field can be replaced with a proper lookup.
	nombreEnfermero := licenseNumber.(string)

	var req application.UpdateEstadoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// nombrePaciente is optionally sent by the client to enrich the event payload.
	// If omitted the event will carry an empty string; the WebSocket consumer can
	// resolve the name from paciente_id if needed.
	nombrePaciente := ctx.GetHeader("X-Paciente-Nombre")

	if err := c.useCase.Execute(
		pacienteID,
		enfermeroID.(string),
		nombreEnfermero,
		nombrePaciente,
		req,
	); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "paciente no encontrado o no asignado a esta enfermera" {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estado actualizado exitosamente"})
}
