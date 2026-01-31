package services

import (
	"fmt"
	"time"

	"github.com/refortunato/go_app_base/internal/shared"
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
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := s.repository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
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
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(name, description string, price float64, stock int) (*models.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if stock < 0 {
		return nil, fmt.Errorf("stock cannot be negative")
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
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id, name, description string, price float64, stock int) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	existing, err := s.repository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("product not found")
	}

	if name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if stock < 0 {
		return nil, fmt.Errorf("stock cannot be negative")
	}

	existing.Name = name
	existing.Description = description
	existing.Price = price
	existing.Stock = stock
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return existing, nil
}

// DeleteProduct removes a product by ID
func (s *ProductService) DeleteProduct(id string) error {
	if id == "" {
		return fmt.Errorf("product ID is required")
	}

	existing, err := s.repository.FindById(id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
