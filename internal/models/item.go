package models

import (
	"time"
)

// Item represents an item in the system
type Item struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents a request to create an item
type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// UpdateItemRequest represents a request to update an item
type UpdateItemRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

// ListItemsResponse represents a paginated list of items
type ListItemsResponse struct {
	Items      []Item `json:"items"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
	Uptime   string            `json:"uptime"`
}
