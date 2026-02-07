package errors

import (
	sharedErrors "github.com/refortunato/go_app_base/internal/shared/errors"
)

var (
	// Product errors
	ErrProductIdRequired = sharedErrors.NewProblemDetails(
		400,
		"Invalid product ID",
		"Product ID is required",
		"SIP1001",
		sharedErrors.ErrorContextBusiness,
	)
	ErrProductNotFound = sharedErrors.NewProblemDetails(
		404,
		"Product not found",
		"The requested product was not found",
		"SIP1002",
		sharedErrors.ErrorContextBusiness,
	)
	ErrProductNameRequired = sharedErrors.NewProblemDetails(
		400,
		"Invalid product name",
		"Product name is required",
		"SIP1003",
		sharedErrors.ErrorContextBusiness,
	)
	ErrProductPriceInvalid = sharedErrors.NewProblemDetails(
		400,
		"Invalid product price",
		"Product price cannot be negative",
		"SIP1004",
		sharedErrors.ErrorContextBusiness,
	)
	ErrProductStockInvalid = sharedErrors.NewProblemDetails(
		400,
		"Invalid product stock",
		"Product stock cannot be negative",
		"SIP1005",
		sharedErrors.ErrorContextBusiness,
	)

	// Generic errors
	ErrGeneric = sharedErrors.NewProblemDetails(
		500,
		"Internal server error",
		"An unexpected error occurred",
		"SIP9999",
		sharedErrors.ErrorContextInfra,
	)
)
