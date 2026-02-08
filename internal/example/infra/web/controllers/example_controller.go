package controllers

import (
	"net/http"

	"github.com/refortunato/go_app_base/internal/example/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/shared/logger"
	"github.com/refortunato/go_app_base/internal/shared/web/advisor"
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

// GetExample godoc
// @Summary      Get example by ID
// @Description  Retrieves a specific example entity from the database
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID format)"
// @Success      200  {object}  usecases.GetExampleOutputDTO
// @Failure      404  {object}  errors.ProblemDetails  "Example not found"
// @Failure      500  {object}  errors.ProblemDetails  "Internal server error"
// @Router       /examples/{id} [get]
func (controller *ExampleController) GetExample(c webcontext.WebContext) {
	id := c.Param("id")
	ctx := c.GetContext()

	// Log the incoming request with custom fields
	logger.Info(ctx, "Processing GetExample request", logger.CustomFields{
		"exampleId": id,
		"endpoint":  "GET /examples/:id",
	})

	input := usecases.GetExampleInputDTO{
		Id: id,
	}

	// Pass context from request for trace propagation
	output, err := controller.GetExampleUseCase.Execute(ctx, input)
	if err != nil {
		// Log error with custom context
		logger.Error(ctx, "Failed to get example", logger.CustomFields{
			"exampleId": id,
			"error":     err.Error(),
		})
		advisor.ReturnApplicationError(c, err)
		return
	}

	// Log successful response
	logger.Info(ctx, "Example retrieved successfully", logger.CustomFields{
		"exampleId": id,
	})

	c.JSON(http.StatusOK, output)
}
