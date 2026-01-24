package usecases

import (
	"github.com/refortunato/go_app_base/internal/health/core/application/repositories"
)

type HealthCheckOutputDTO struct {
	Status string `json:"status"`
}

type HealthCheckUseCase struct {
	healthRepository repositories.HealthRepository
}

func NewHealthCheckUseCase(healthRepository repositories.HealthRepository) *HealthCheckUseCase {
	return &HealthCheckUseCase{
		healthRepository: healthRepository,
	}
}

func (u *HealthCheckUseCase) Execute() (*HealthCheckOutputDTO, error) {
	err := u.healthRepository.CheckDatabaseConnection()
	if err != nil {
		return nil, err
	}

	output := &HealthCheckOutputDTO{
		Status: "OK",
	}

	return output, nil
}
