package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/admin/application"

	"github.com/gin-gonic/gin"
)

// GetAreasController maneja la solicitud de distribución de pacientes por área (tab Áreas).
type GetAreasController struct {
	useCase *application.GetAreasUseCase
}

func NewGetAreasController(uc *application.GetAreasUseCase) *GetAreasController {
	return &GetAreasController{useCase: uc}
}

// GetAreas godoc
// @Summary      Distribución por área hospitalaria
// @Description  Devuelve la cantidad de pacientes (total y críticos) por cada área del hospital.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   entities.AreaDistribucion
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/admin/areas [get]
func (c *GetAreasController) GetAreas(ctx *gin.Context) {
	areas, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, areas)
}
