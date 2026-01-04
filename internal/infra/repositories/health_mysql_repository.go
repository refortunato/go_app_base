package repositories

import (
	"database/sql"
)

type HealthMySQLRepository struct {
	db *sql.DB
}

func NewHealthMySQLRepository(db *sql.DB) *HealthMySQLRepository {
	return &HealthMySQLRepository{db: db}
}

func (r *HealthMySQLRepository) CheckDatabaseConnection() error {
	// Simple query to check database connectivity
	var result int
	err := r.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return err
	}
	return nil
}
