package advisor

import (
	"net/http"

	app_errors "github.com/refortunato/go_app_base/internal/shared/errors"
	webcontext "github.com/refortunato/go_app_base/internal/shared/web/context"
)

func ReturnApplicationError(c webcontext.WebContext, err error) {
	if err != nil {
		// Retornar erros formatados como ProblemDetails
		if pd, ok := err.(*app_errors.ProblemDetails); ok {
			c.JSON(pd.Status, pd)
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not execute operation"})
		return
	}
}

func ReturnBadRequestError(c webcontext.WebContext, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
}

func ReturnNotFoundError(c webcontext.WebContext) {
	c.JSON(http.StatusNotFound, map[string]string{"error": "resource not found"})
}
