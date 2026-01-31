package dto

import (
	"errors"
	"strconv"
)

// PaginationRequestDTO represents pagination parameters for list queries
type PaginationRequestDTO struct {
	Page   int
	Limit  int
	Offset int
}

// NewPaginationRequestDTO creates a pagination DTO from query string parameters
// Default values: page=1, limit=10
func NewPaginationRequestDTO(pageStr, limitStr string) (*PaginationRequestDTO, error) {
	page := 1
	limit := 10

	// Parse page
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			return nil, errors.New("invalid page parameter")
		}
	}

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		} else {
			return nil, errors.New("invalid limit parameter")
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	return &PaginationRequestDTO{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// PaginationResponseDTO represents pagination metadata in responses
type PaginationResponseDTO struct {
	Data       int `json:"data"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// NewPaginationResponseDTO creates pagination metadata for responses
func NewPaginationResponseDTO(data, limit, totalItems int) *PaginationResponseDTO {
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return &PaginationResponseDTO{
		Data:       data,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
