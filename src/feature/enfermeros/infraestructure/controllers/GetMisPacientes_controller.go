package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/enfermeros/application"

	"github.com/gin-gonic/gin"
)

type GetMisPacientesController struct {
	useCase *application.GetMisPacientesUseCase
}

func NewGetMisPacientesController(uc *application.GetMisPacientesUseCase) *GetMisPacientesController {
	return &GetMisPacientesController{useCase: uc}
}

// GetMisPacientes godoc
// @Summary      Pacientes asignados a la enfermera
// @Description  Retorna la lista de pacientes asignados a la enfermera autenticada
//               junto con los contadores de estado (Total, Críticos, Observación, Estables).
// @Tags         enfermeros
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  entities.MisPacientesResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/enfermeros/mis-pacientes [get]
func (c *GetMisPacientesController) GetMisPacientes(ctx *gin.Context) {
	// The nurse's user_id is injected by AuthMiddleware from the JWT.
	enfermeroID, _ := ctx.Get("user_id")

	response, err := c.useCase.Execute(enfermeroID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
