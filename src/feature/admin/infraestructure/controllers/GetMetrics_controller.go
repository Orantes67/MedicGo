package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/admin/application"

	"github.com/gin-gonic/gin"
)

// GetMetricsController maneja la solicitud de métricas del panel administrativo.
type GetMetricsController struct {
	useCase *application.GetMetricsUseCase
}

func NewGetMetricsController(uc *application.GetMetricsUseCase) *GetMetricsController {
	return &GetMetricsController{useCase: uc}
}

// GetMetrics godoc
// @Summary      Métricas del panel administrativo
// @Description  Devuelve totales de pacientes por estado, personal activo, alertas y actividad reciente.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  entities.AdminMetrics
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/admin/metricas [get]
func (c *GetMetricsController) GetMetrics(ctx *gin.Context) {
	metrics, err := c.useCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}
