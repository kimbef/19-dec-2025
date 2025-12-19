package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kimbef/19-dec-2025/internal/models"
	"github.com/kimbef/19-dec-2025/internal/repository"
	"github.com/kimbef/19-dec-2025/pkg/logger"
	"github.com/kimbef/19-dec-2025/pkg/response"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	itemRepo  *repository.ItemRepository
	log       *logger.Logger
	startTime time.Time
}

// New creates a new handler instance
func New(itemRepo *repository.ItemRepository, log *logger.Logger) *Handler {
	return &Handler{
		itemRepo:  itemRepo,
		log:       log,
		startTime: time.Now(),
	}
}

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	healthResp := models.HealthResponse{
		Status: "healthy",
		Services: map[string]string{
			"api":      "up",
			"database": "up",
		},
		Uptime: uptime.String(),
	}

	response.Success(w, healthResp)
}

// CreateItem handles POST /api/items
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req models.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Failed to decode request: %v", err)
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if req.Name == "" {
		response.BadRequest(w, "Name is required")
		return
	}
	if req.Quantity < 0 {
		response.BadRequest(w, "Quantity must be non-negative")
		return
	}
	if req.Price < 0 {
		response.BadRequest(w, "Price must be non-negative")
		return
	}

	item, err := h.itemRepo.Create(r.Context(), req)
	if err != nil {
		h.log.Errorf("Failed to create item: %v", err)
		response.InternalError(w, "Failed to create item")
		return
	}

	response.Created(w, item)
}

// GetItem handles GET /api/items/{id}
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid item ID")
		return
	}

	item, err := h.itemRepo.GetByID(r.Context(), id)
	if err != nil {
		h.log.Errorf("Failed to get item: %v", err)
		response.InternalError(w, "Failed to retrieve item")
		return
	}

	if item == nil {
		response.NotFound(w, "Item not found")
		return
	}

	response.Success(w, item)
}

// UpdateItem handles PUT /api/items/{id}
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid item ID")
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warnf("Failed to decode request: %v", err)
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if req.Quantity != nil && *req.Quantity < 0 {
		response.BadRequest(w, "Quantity must be non-negative")
		return
	}
	if req.Price != nil && *req.Price < 0 {
		response.BadRequest(w, "Price must be non-negative")
		return
	}

	item, err := h.itemRepo.Update(r.Context(), id, req)
	if err != nil {
		h.log.Errorf("Failed to update item: %v", err)
		response.InternalError(w, "Failed to update item")
		return
	}

	if item == nil {
		response.NotFound(w, "Item not found")
		return
	}

	response.Success(w, item)
}

// DeleteItem handles DELETE /api/items/{id}
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid item ID")
		return
	}

	err = h.itemRepo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		response.NotFound(w, "Item not found")
		return
	}
	if err != nil {
		h.log.Errorf("Failed to delete item: %v", err)
		response.InternalError(w, "Failed to delete item")
		return
	}

	response.NoContent(w)
}

// ListItems handles GET /api/items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	perPage := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	items, total, err := h.itemRepo.List(r.Context(), page, perPage)
	if err != nil {
		h.log.Errorf("Failed to list items: %v", err)
		response.InternalError(w, "Failed to retrieve items")
		return
	}

	// Calculate total pages
	totalPages := (total + perPage - 1) / perPage

	resp := models.ListItemsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	response.Success(w, resp)
}
