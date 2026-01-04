package errors

var (
	ErrDescriptionIsRequired = NewProblemDetails(
		400,
		"Invalid description",
		"Description is required and cannot be empty",
		"EX1001",
		ErrorContextBusiness,
	)
	ErrExampleNotFound = NewProblemDetails(
		404,
		"Example not found",
		"The requested example was not found",
		"EX1002",
		ErrorContextBusiness,
	)
	ErrDatabaseConnection = NewProblemDetails(
		500,
		"Database connection error",
		"Failed to connect to the database",
		"DB1001",
		ErrorContextInfra,
	)
)

const (
	ErrorContextBusiness = "business"
	ErrorContextInfra    = "infra"
	ErrorContextGeneric  = "generic"
)
