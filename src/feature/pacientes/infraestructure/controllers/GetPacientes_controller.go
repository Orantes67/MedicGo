package controllers

import (
	"net/http"

	"ApiMedicGO/src/feature/pacientes/application"

	"github.com/gin-gonic/gin"
)

type GetPacientesController struct {
	useCase *application.GetPacientesUseCase
}

func NewGetPacientesController(uc *application.GetPacientesUseCase) *GetPacientesController {
	return &GetPacientesController{useCase: uc}
}

func (c *GetPacientesController) GetAll(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")

	pacientes, err := c.useCase.Execute(role.(string), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":     len(pacientes),
		"pacientes": pacientes,
	})
}
