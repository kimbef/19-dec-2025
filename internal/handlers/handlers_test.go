package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kimbef/19-dec-2025/internal/models"
	"github.com/kimbef/19-dec-2025/pkg/logger"
)

func TestHealth(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp["success"].(bool) {
		t.Error("expected success to be true")
	}

	data := resp["data"].(map[string]interface{})
	if data["status"] != "healthy" {
		t.Errorf("expected status healthy, got %v", data["status"])
	}
}

func TestCreateItem_InvalidJSON(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateItem_MissingName(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	reqBody := models.CreateItemRequest{
		Description: "Test Description",
		Quantity:    10,
		Price:       99.99,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateItem_NegativeQuantity(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	reqBody := models.CreateItemRequest{
		Name:        "Test Item",
		Description: "Test Description",
		Quantity:    -5,
		Price:       99.99,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateItem_NegativePrice(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	reqBody := models.CreateItemRequest{
		Name:        "Test Item",
		Description: "Test Description",
		Quantity:    10,
		Price:       -99.99,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetItem_InvalidID(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/items/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	h.GetItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateItem_InvalidID(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	reqBody := models.UpdateItemRequest{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/items/invalid", bytes.NewBuffer(body))
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	h.UpdateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteItem_InvalidID(t *testing.T) {
	log := logger.NewDefault()
	h := &Handler{
		log:       log,
		startTime: time.Now(),
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/items/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	h.DeleteItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

