package entities

import (
	"time"

	"github.com/refortunato/go_app_base/internal/shared"
	sharedErrors "github.com/refortunato/go_app_base/internal/shared/errors"
)

type Example struct {
	id          string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewExample(description string) (*Example, error) {
	example := &Example{
		id:          shared.GenerateId(),
		description: description,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}
	if err := example.Validate(); err != nil {
		return nil, err
	}
	return example, nil
}

func RestoreExample(
	id,
	description string,
	createdAt,
	updatedAt time.Time) (*Example, error) {
	return &Example{
		id:          id,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (e *Example) Validate() error {
	if e.description == "" {
		return sharedErrors.ErrDescriptionIsRequired
	}
	return nil
}

// Getters

func (e *Example) GetId() string {
	return e.id
}

func (e *Example) GetDescription() string {
	return e.description
}

func (e *Example) GetCreatedAt() time.Time {
	return e.createdAt
}

func (e *Example) GetUpdatedAt() time.Time {
	return e.updatedAt
}

// Setters

func (e *Example) SetDescription(description string) {
	e.description = description
	e.updatedAt = time.Now().UTC()
}
