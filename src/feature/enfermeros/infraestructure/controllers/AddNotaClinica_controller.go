package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/enfermeros/application"

	"github.com/gin-gonic/gin"
)

type AddNotaClinicaEnfermeroController struct {
	useCase *application.AddNotaClinicaEnfermeroUseCase
}

func NewAddNotaClinicaEnfermeroController(uc *application.AddNotaClinicaEnfermeroUseCase) *AddNotaClinicaEnfermeroController {
	return &AddNotaClinicaEnfermeroController{useCase: uc}
}

// AddNota godoc
// @Summary      Agregar nota clínica (enfermera)
// @Description  La enfermera escribe una nueva nota clínica para el paciente asignado.
//               Dispara un PacienteEstadoEvent "nota_clinica" listo para el WebSocket.
// @Tags         enfermeros
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                                   true  "ID del paciente"
// @Param        body  body  application.AddNotaEnfermeroRequest      true  "Contenido de la nota"
// @Success      201  {object}  entities.PatientNote
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/enfermeros/pacientes/{id}/notas [post]
func (c *AddNotaClinicaEnfermeroController) AddNota(ctx *gin.Context) {
	pacienteID := ctx.Param("id")
	enfermeroID, _   := ctx.Get("user_id")
	licenseNumber, _ := ctx.Get("license_number")
	nombrePaciente   := ctx.GetHeader("X-Paciente-Nombre")

	var req application.AddNotaEnfermeroRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nota, err := c.useCase.Execute(
		pacienteID,
		enfermeroID.(string),
		licenseNumber.(string),
		nombrePaciente,
		req,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, nota)
}
