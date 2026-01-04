package repositories

import (
	"database/sql"
	"time"

	"github.com/refortunato/go_app_base/internal/core/domain/entities"
)

type exampleEntity struct {
	Id          string    `db:"id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ExampleMySQLRepository struct {
	db *sql.DB
}

func NewExampleMySQLRepository(db *sql.DB) *ExampleMySQLRepository {
	return &ExampleMySQLRepository{db: db}
}

func (r *ExampleMySQLRepository) Save(example *entities.Example) error {
	stmt, err := r.db.Prepare("INSERT INTO examples (id, description, created_at, updated_at) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		example.GetId(),
		example.GetDescription(),
		example.GetCreatedAt(),
		example.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExampleMySQLRepository) FindById(id string) (*entities.Example, error) {
	row := r.db.QueryRow("SELECT id, description, created_at, updated_at FROM examples WHERE id = ?", id)
	var exampleEntity exampleEntity
	err := row.Scan(
		&exampleEntity.Id,
		&exampleEntity.Description,
		&exampleEntity.CreatedAt,
		&exampleEntity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	exampleDomain, err := r.mapToDomain(exampleEntity)
	if err != nil {
		return nil, err
	}
	return exampleDomain, nil
}

func (r *ExampleMySQLRepository) Update(example *entities.Example) error {
	stmt, err := r.db.Prepare("UPDATE examples SET description=?, updated_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		example.GetDescription(),
		example.GetUpdatedAt(),
		example.GetId(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExampleMySQLRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM examples WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExampleMySQLRepository) mapToDomain(entity exampleEntity) (*entities.Example, error) {
	return entities.RestoreExample(
		entity.Id,
		entity.Description,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}
