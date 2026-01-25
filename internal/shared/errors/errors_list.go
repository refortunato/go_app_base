package errors

// Generic infrastructure errors
var (
	ErrDatabaseConnection = NewProblemDetails(
		500,
		"Database connection error",
		"Failed to connect to the database",
		"DB1001",
		ErrorContextInfra,
	)
)
