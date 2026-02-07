package services

import (
	"time"

	"github.com/refortunato/go_app_base/internal/shared"
	"github.com/refortunato/go_app_base/internal/simple_module/errors"
	"github.com/refortunato/go_app_base/internal/simple_module/models"
	"github.com/refortunato/go_app_base/internal/simple_module/repositories"
)

// ProductService handles business logic for products
type ProductService struct {
	repository *repositories.ProductRepository
}

// NewProductService creates a new product service instance
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repository: repo}
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	if id == "" {
		return nil, errors.ErrProductIdRequired
	}

	product, err := s.repository.FindById(id)
	if err != nil {
		return nil, errors.ErrGeneric
	}

	if product == nil {
		return nil, errors.ErrProductNotFound
	}

	return product, nil
}

// ListProducts retrieves all products with pagination
func (s *ProductService) ListProducts(limit, offset int) ([]*models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repository.FindAll(limit, offset)
	if err != nil {
		return nil, errors.ErrGeneric
	}

	return products, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(name, description string, price float64, stock int) (*models.Product, error) {
	if name == "" {
		return nil, errors.ErrProductNameRequired
	}
	if price < 0 {
		return nil, errors.ErrProductPriceInvalid
	}
	if stock < 0 {
		return nil, errors.ErrProductStockInvalid
	}

	now := time.Now().UTC()
	product := &models.Product{
		ID:          shared.GenerateId(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Save(product); err != nil {
		return nil, errors.ErrGeneric
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id, name, description string, price float64, stock int) (*models.Product, error) {
	if id == "" {
		return nil, errors.ErrProductIdRequired
	}

	existing, err := s.repository.FindById(id)
	if err != nil {
		return nil, errors.ErrGeneric
	}
	if existing == nil {
		return nil, errors.ErrProductNotFound
	}

	if name == "" {
		return nil, errors.ErrProductNameRequired
	}
	if price < 0 {
		return nil, errors.ErrProductPriceInvalid
	}
	if stock < 0 {
		return nil, errors.ErrProductStockInvalid
	}

	existing.Name = name
	existing.Description = description
	existing.Price = price
	existing.Stock = stock
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(existing); err != nil {
		return nil, errors.ErrGeneric
	}

	return existing, nil
}

// DeleteProduct removes a product by ID
func (s *ProductService) DeleteProduct(id string) error {
	if id == "" {
		return errors.ErrProductIdRequired
	}

	existing, err := s.repository.FindById(id)
	if err != nil {
		return errors.ErrGeneric
	}
	if existing == nil {
		return errors.ErrProductNotFound
	}

	if err := s.repository.Delete(id); err != nil {
		return errors.ErrGeneric
	}

	return nil
}
