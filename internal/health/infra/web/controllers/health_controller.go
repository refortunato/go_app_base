package controllers

import (
	"net/http"

	"github.com/refortunato/go_app_base/internal/health/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/shared/web/advisor"
	webcontext "github.com/refortunato/go_app_base/internal/shared/web/context"
)

type HealthController struct {
	HealthCheckUseCase usecases.HealthCheckUseCase
}

func NewHealthController(healthCheckUseCase usecases.HealthCheckUseCase) *HealthController {
	return &HealthController{
		HealthCheckUseCase: healthCheckUseCase,
	}
}

func (controller *HealthController) HealthCheck(c webcontext.WebContext) {
	output, err := controller.HealthCheckUseCase.Execute()
	if err != nil {
		advisor.ReturnApplicationError(c, err)
		return
	}
	c.JSON(http.StatusOK, output)
}
