package models

import "time"

// Product represents a simple product data structure
type Product struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Laptop Dell XPS 15"`
	Description string    `json:"description" example:"High-performance laptop for professionals"`
	Price       float64   `json:"price" example:"5499.99"`
	Stock       int       `json:"stock" example:"10"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T10:00:00Z"`
}
