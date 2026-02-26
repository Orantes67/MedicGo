package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/doctores/application"

	"github.com/gin-gonic/gin"
)

type AddNotaClinicaController struct {
	useCase *application.AddNotaClinicaUseCase
}

func NewAddNotaClinicaController(uc *application.AddNotaClinicaUseCase) *AddNotaClinicaController {
	return &AddNotaClinicaController{useCase: uc}
}

// AddNotaClinica godoc
// @Summary      Agregar nota clínica
// @Description  El doctor escribe una nueva nota para el paciente.
//               Dispara un DoctorEvent "nota_clinica" listo para el WebSocket.
// @Tags         doctores
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                          true  "ID del paciente"
// @Param        body  body  application.AddNotaRequest      true  "Contenido de la nota"
// @Success      201  {object}  entities.PatientNote
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/doctores/pacientes/{id}/notas [post]
func (c *AddNotaClinicaController) AddNota(ctx *gin.Context) {
	pacienteID := ctx.Param("id")
	doctorID, _      := ctx.Get("user_id")
	licenseNumber, _ := ctx.Get("license_number")

	var req application.AddNotaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nombrePaciente := ctx.GetHeader("X-Paciente-Nombre")

	nota, err := c.useCase.Execute(
		pacienteID,
		doctorID.(string),
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
