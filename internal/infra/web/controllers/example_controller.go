package controllers

import (
	"net/http"

	"github.com/refortunato/go_app_base/internal/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/infra/web/controllers/advisor"
	"github.com/refortunato/go_app_base/internal/infra/web/webcontext"
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

	input := usecases.GetExampleInputDTO{
		Id: id,
	}

	output, err := controller.GetExampleUseCase.Execute(input)
	if err != nil {
		advisor.ReturnApplicationError(c, err)
		return
	}
	c.JSON(http.StatusOK, output)
}
