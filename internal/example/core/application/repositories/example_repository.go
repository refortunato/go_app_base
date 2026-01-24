package repositories

import (
	"github.com/refortunato/go_app_base/internal/example/core/domain/entities"
)

type ExampleRepository interface {
	Save(example *entities.Example) error
	FindById(id string) (*entities.Example, error)
	Update(example *entities.Example) error
	Delete(id string) error
}
