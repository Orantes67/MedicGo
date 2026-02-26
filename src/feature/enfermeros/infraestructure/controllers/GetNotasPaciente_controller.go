package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/enfermeros/application"

	"github.com/gin-gonic/gin"
)

type GetNotasPacienteController struct {
	useCase *application.GetNotasPacienteUseCase
}

func NewGetNotasPacienteController(uc *application.GetNotasPacienteUseCase) *GetNotasPacienteController {
	return &GetNotasPacienteController{useCase: uc}
}

// GetNotas godoc
// @Summary      Notas clínicas del paciente (vista enfermera)
// @Description  Retorna todas las notas clínicas del paciente (escritas por doctor o enfermera).
//               Solo disponible si el paciente está asignado a la enfermera autenticada.
// @Tags         enfermeros
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "ID del paciente"
// @Success      200  {array}   entities.PatientNote
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/enfermeros/pacientes/{id}/notas [get]
func (c *GetNotasPacienteController) GetNotas(ctx *gin.Context) {
	pacienteID := ctx.Param("id")
	enfermeroID, _ := ctx.Get("user_id")

	notas, err := c.useCase.Execute(pacienteID, enfermeroID.(string))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "paciente no encontrado o no asignado a esta enfermera" {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notas)
}
