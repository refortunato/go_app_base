package controllers

import (
	"net/http"

	"github.com/refortunato/go_app_base/internal/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/infra/web/controllers/advisor"
	"github.com/refortunato/go_app_base/internal/shared/logger"
	webcontext "github.com/refortunato/go_app_base/internal/shared/web/context"
)

type ExampleController struct {
	GetExampleUseCase usecases.GetExampleUseCase
}

func NewExampleController(getExampleUseCase usecases.GetExampleUseCase) *ExampleController {
	return &ExampleController{
		GetExampleUseCase: getExampleUseCase,
	}
}

func (controller *ExampleController) GetExample(c webcontext.WebContext) {
	id := c.Param("id")

	// Log the incoming request with custom fields
	logger.Info("Processing GetExample request", logger.CustomFields{
		"exampleId": id,
		"endpoint":  "GET /examples/:id",
	})

	input := usecases.GetExampleInputDTO{
		Id: id,
	}

	output, err := controller.GetExampleUseCase.Execute(input)
	if err != nil {
		// Log error with custom context
		logger.Error("Failed to get example", logger.CustomFields{
			"exampleId": id,
			"error":     err.Error(),
		})
		advisor.ReturnApplicationError(c, err)
		return
	}

	// Log successful response
	logger.Info("Example retrieved successfully", logger.CustomFields{
		"exampleId": id,
	})

	c.JSON(http.StatusOK, output)
}
