package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kimbef/19-dec-2025/internal/models"
	"github.com/kimbef/19-dec-2025/pkg/logger"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	t.Run("health check without database", func(t *testing.T) {
		handler := &Handler{
			log: log,
			db:  nil,
		}

		router := gin.New()
		router.GET("/health", handler.HealthCheck)

		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		var health models.HealthCheck
		if err := json.Unmarshal(w.Body.Bytes(), &health); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if health.Status != "unhealthy" {
			t.Errorf("Expected status unhealthy, got %s", health.Status)
		}
	})
}

func TestReadyCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	t.Run("ready check without database", func(t *testing.T) {
		handler := &Handler{
			log: log,
			db:  nil,
		}

		router := gin.New()
		router.GET("/ready", handler.ReadyCheck)

		req, _ := http.NewRequest("GET", "/ready", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if ready, ok := response["ready"].(bool); !ok || ready {
			t.Errorf("Expected ready to be false")
		}
	})
}

func TestCreatePipelineValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	handler := &Handler{
		log: log,
		db:  nil,
	}

	router := gin.New()
	router.POST("/api/v1/pipelines", handler.CreatePipeline)

	t.Run("invalid request body - missing required field", func(t *testing.T) {
		reqBody := []byte(`{"description": "Missing name field"}`)
		req, _ := http.NewRequest("POST", "/api/v1/pipelines", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetPipelineInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	handler := &Handler{
		log: log,
		db:  nil,
	}

	router := gin.New()
	router.GET("/api/v1/pipelines/:id", handler.GetPipeline)

	req, _ := http.NewRequest("GET", "/api/v1/pipelines/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Errorf("Expected error in response")
	}
}

func TestCreateBuildValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	handler := &Handler{
		log: log,
		db:  nil,
	}

	router := gin.New()
	router.POST("/api/v1/builds", handler.CreateBuild)

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := []byte(`{"branch": "main"}`)
		req, _ := http.NewRequest("POST", "/api/v1/builds", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestCreateDeploymentValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.GetDefaultLogger()

	handler := &Handler{
		log: log,
		db:  nil,
	}

	router := gin.New()
	router.POST("/api/v1/deployments", handler.CreateDeployment)

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := []byte(`{"environment": "production"}`)
		req, _ := http.NewRequest("POST", "/api/v1/deployments", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
