package usecases

import (
	"time"

	"github.com/refortunato/go_app_base/internal/core/application/repositories"
)

type GetExampleInputDTO struct {
	Id string
}

type GetExampleOutputDTO struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetExampleUseCase struct {
	exampleRepository repositories.ExampleRepository
}

func NewGetExampleUseCase(exampleRepository repositories.ExampleRepository) *GetExampleUseCase {
	return &GetExampleUseCase{
		exampleRepository: exampleRepository,
	}
}

func (u *GetExampleUseCase) Execute(input GetExampleInputDTO) (*GetExampleOutputDTO, error) {
	example, err := u.exampleRepository.FindById(input.Id)
	if err != nil {
		return nil, err
	}

	output := &GetExampleOutputDTO{
		Id:          example.GetId(),
		Description: example.GetDescription(),
		CreatedAt:   example.GetCreatedAt(),
		UpdatedAt:   example.GetUpdatedAt(),
	}

	return output, nil
}
