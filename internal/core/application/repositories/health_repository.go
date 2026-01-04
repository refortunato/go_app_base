package repositories

type HealthRepository interface {
	CheckDatabaseConnection() error
}
