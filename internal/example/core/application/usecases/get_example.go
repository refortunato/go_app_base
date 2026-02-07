package usecases

import (
	"context"
	"time"

	"github.com/refortunato/go_app_base/internal/example/core/application/repositories"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type GetExampleInputDTO struct {
	Id string
}

type GetExampleOutputDTO struct {
	Id          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description string    `json:"description" example:"Sample example description"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T10:00:00Z"`
}

type GetExampleUseCase struct {
	exampleRepository repositories.ExampleRepository
}

func NewGetExampleUseCase(exampleRepository repositories.ExampleRepository) *GetExampleUseCase {
	return &GetExampleUseCase{
		exampleRepository: exampleRepository,
	}
}

func (u *GetExampleUseCase) Execute(ctx context.Context, input GetExampleInputDTO) (*GetExampleOutputDTO, error) {
	// Create a span for this use case execution
	tracer := otel.Tracer("example.usecase")
	ctx, span := tracer.Start(ctx, "GetExampleUseCase.Execute")
	defer span.End()

	// Add attributes to the span for better observability
	span.SetAttributes(
		attribute.String("example.id", input.Id),
		attribute.String("usecase", "GetExample"),
	)

	example, err := u.exampleRepository.FindById(input.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find example")
		return nil, err
	}

	output := &GetExampleOutputDTO{
		Id:          example.GetId(),
		Description: example.GetDescription(),
		CreatedAt:   example.GetCreatedAt(),
		UpdatedAt:   example.GetUpdatedAt(),
	}

	span.SetStatus(codes.Ok, "Example retrieved successfully")
	return output, nil
}
