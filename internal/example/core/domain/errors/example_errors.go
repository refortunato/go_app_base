package errors

import (
	sharedErrors "github.com/refortunato/go_app_base/internal/shared/errors"
)

var (
	ErrDescriptionIsRequired = sharedErrors.NewProblemDetails(
		400,
		"Invalid description",
		"Description is required and cannot be empty",
		"EX1001",
		sharedErrors.ErrorContextBusiness,
	)
	ErrExampleNotFound = sharedErrors.NewProblemDetails(
		404,
		"Example not found",
		"The requested example was not found",
		"EX1002",
		sharedErrors.ErrorContextBusiness,
	)
)
