package models

import "time"

// Product represents a simple product data structure
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
