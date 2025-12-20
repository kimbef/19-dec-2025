package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kimbef/19-dec-2025/internal/models"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create creates a new item
func (r *ItemRepository) Create(ctx context.Context, req models.CreateItemRequest) (*models.Item, error) {
	query := `
		INSERT INTO items (name, description, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, quantity, price, created_at, updated_at
	`

	now := time.Now()
	var item models.Item

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Name,
		req.Description,
		req.Quantity,
		req.Price,
		now,
		now,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &item, nil
}

// GetByID retrieves an item by ID
func (r *ItemRepository) GetByID(ctx context.Context, id int) (*models.Item, error) {
	query := `
		SELECT id, name, description, quantity, price, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	var item models.Item
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

// Update updates an existing item
func (r *ItemRepository) Update(ctx context.Context, id int, req models.UpdateItemRequest) (*models.Item, error) {
	// Build dynamic update query
	query := "UPDATE items SET updated_at = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, *req.Name)
		argCount++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argCount)
		args = append(args, *req.Description)
		argCount++
	}
	if req.Quantity != nil {
		query += fmt.Sprintf(", quantity = $%d", argCount)
		args = append(args, *req.Quantity)
		argCount++
	}
	if req.Price != nil {
		query += fmt.Sprintf(", price = $%d", argCount)
		args = append(args, *req.Price)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, description, quantity, price, created_at, updated_at", argCount)
	args = append(args, id)

	var item models.Item
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return &item, nil
}

// Delete deletes an item by ID
func (r *ItemRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM items WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// List retrieves a paginated list of items
func (r *ItemRepository) List(ctx context.Context, page, perPage int) ([]models.Item, int, error) {
	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM items"
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count items: %w", err)
	}

	// Get paginated items
	offset := (page - 1) * perPage
	query := `
		SELECT id, name, description, quantity, price, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	items := make([]models.Item, 0, perPage)
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Quantity,
			&item.Price,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating items: %w", err)
	}

	return items, total, nil
}
