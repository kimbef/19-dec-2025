package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimbef/19-dec-2025/internal/database"
	"github.com/kimbef/19-dec-2025/internal/models"
	"github.com/kimbef/19-dec-2025/pkg/logger"
	"go.uber.org/zap"
)

// Handler contains dependencies for API handlers
type Handler struct {
	db  *database.DB
	log *logger.Logger
}

// NewHandler creates a new API handler
func NewHandler(db *database.DB, log *logger.Logger) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthCheck
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	checks := make(map[string]string)

	// Check database
	if h.db == nil {
		checks["database"] = "unhealthy: database not initialized"
	} else if err := h.db.Health(); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
	} else {
		checks["database"] = "healthy"
	}

	status := "healthy"
	for _, check := range checks {
		if check != "healthy" {
			status = "unhealthy"
			break
		}
	}

	health := models.HealthCheck{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// ReadyCheck godoc
// @Summary Readiness check endpoint
// @Description Returns whether the application is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *Handler) ReadyCheck(c *gin.Context) {
	// Check if database is accessible
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":   false,
			"message": "database not initialized",
		})
		return
	}

	if err := h.db.Health(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":   false,
			"message": "database not ready",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready":   true,
		"message": "application is ready",
	})
}

// ListPipelines godoc
// @Summary List all pipelines
// @Description Get a list of all CI/CD pipelines
// @Tags pipelines
// @Produce json
// @Success 200 {array} models.Pipeline
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines [get]
func (h *Handler) ListPipelines(c *gin.Context) {
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at
		FROM pipelines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		h.log.Error("Failed to query pipelines", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipelines"})
		return
	}
	defer rows.Close()

	var pipelines []models.Pipeline
	for rows.Next() {
		var p models.Pipeline
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			h.log.Error("Failed to scan pipeline", zap.Error(err))
			continue
		}
		pipelines = append(pipelines, p)
	}

	if pipelines == nil {
		pipelines = []models.Pipeline{}
	}

	c.JSON(http.StatusOK, pipelines)
}

// GetPipeline godoc
// @Summary Get a pipeline by ID
// @Description Get details of a specific pipeline
// @Tags pipelines
// @Produce json
// @Param id path int true "Pipeline ID"
// @Success 200 {object} models.Pipeline
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines/{id} [get]
func (h *Handler) GetPipeline(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at
		FROM pipelines
		WHERE id = $1 AND deleted_at IS NULL
	`

	var p models.Pipeline
	err = h.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}
	if err != nil {
		h.log.Error("Failed to query pipeline", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// CreatePipeline godoc
// @Summary Create a new pipeline
// @Description Create a new CI/CD pipeline
// @Tags pipelines
// @Accept json
// @Produce json
// @Param pipeline body models.CreatePipelineRequest true "Pipeline data"
// @Success 201 {object} models.Pipeline
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines [post]
func (h *Handler) CreatePipeline(c *gin.Context) {
	var req models.CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO pipelines (name, description, status)
		VALUES ($1, $2, 'active')
		RETURNING id, name, description, status, created_at, updated_at, deleted_at
	`

	var p models.Pipeline
	err := h.db.QueryRow(query, req.Name, req.Description).Scan(
		&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		h.log.Error("Failed to create pipeline", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// UpdatePipeline godoc
// @Summary Update a pipeline
// @Description Update an existing pipeline
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path int true "Pipeline ID"
// @Param pipeline body models.UpdatePipelineRequest true "Pipeline data"
// @Success 200 {object} models.Pipeline
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines/{id} [put]
func (h *Handler) UpdatePipeline(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req models.UpdatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE pipelines
		SET name = COALESCE(NULLIF($1, ''), name),
		    description = COALESCE(NULLIF($2, ''), description),
		    status = COALESCE(NULLIF($3, ''), status),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, name, description, status, created_at, updated_at, deleted_at
	`

	var p models.Pipeline
	err = h.db.QueryRow(query, req.Name, req.Description, req.Status, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}
	if err != nil {
		h.log.Error("Failed to update pipeline", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pipeline"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// DeletePipeline godoc
// @Summary Delete a pipeline
// @Description Soft delete a pipeline
// @Tags pipelines
// @Produce json
// @Param id path int true "Pipeline ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines/{id} [delete]
func (h *Handler) DeletePipeline(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	query := `
		UPDATE pipelines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := h.db.Exec(query, id)
	if err != nil {
		h.log.Error("Failed to delete pipeline", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pipeline"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListBuilds godoc
// @Summary List builds for a pipeline
// @Description Get a list of all builds for a specific pipeline
// @Tags builds
// @Produce json
// @Param pipeline_id path int true "Pipeline ID"
// @Success 200 {array} models.Build
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pipelines/{pipeline_id}/builds [get]
func (h *Handler) ListBuilds(c *gin.Context) {
	pipelineID, err := strconv.ParseInt(c.Param("pipeline_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	query := `
		SELECT id, pipeline_id, build_number, status, branch, commit_hash,
		       started_at, finished_at, duration, created_at, updated_at
		FROM builds
		WHERE pipeline_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query, pipelineID)
	if err != nil {
		h.log.Error("Failed to query builds", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch builds"})
		return
	}
	defer rows.Close()

	var builds []models.Build
	for rows.Next() {
		var b models.Build
		if err := rows.Scan(&b.ID, &b.PipelineID, &b.BuildNumber, &b.Status, &b.Branch, &b.CommitHash,
			&b.StartedAt, &b.FinishedAt, &b.Duration, &b.CreatedAt, &b.UpdatedAt); err != nil {
			h.log.Error("Failed to scan build", zap.Error(err))
			continue
		}
		builds = append(builds, b)
	}

	if builds == nil {
		builds = []models.Build{}
	}

	c.JSON(http.StatusOK, builds)
}

// CreateBuild godoc
// @Summary Create a new build
// @Description Create a new build for a pipeline
// @Tags builds
// @Accept json
// @Produce json
// @Param build body models.CreateBuildRequest true "Build data"
// @Success 201 {object} models.Build
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/builds [post]
func (h *Handler) CreateBuild(c *gin.Context) {
	var req models.CreateBuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get next build number for this pipeline
	var buildNumber int
	err := h.db.QueryRow(`
		SELECT COALESCE(MAX(build_number), 0) + 1
		FROM builds
		WHERE pipeline_id = $1
	`, req.PipelineID).Scan(&buildNumber)
	if err != nil {
		h.log.Error("Failed to get next build number", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create build"})
		return
	}

	query := `
		INSERT INTO builds (pipeline_id, build_number, status, branch, commit_hash)
		VALUES ($1, $2, 'pending', $3, $4)
		RETURNING id, pipeline_id, build_number, status, branch, commit_hash,
		          started_at, finished_at, duration, created_at, updated_at
	`

	var b models.Build
	err = h.db.QueryRow(query, req.PipelineID, buildNumber, req.Branch, req.CommitHash).Scan(
		&b.ID, &b.PipelineID, &b.BuildNumber, &b.Status, &b.Branch, &b.CommitHash,
		&b.StartedAt, &b.FinishedAt, &b.Duration, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		h.log.Error("Failed to create build", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create build"})
		return
	}

	c.JSON(http.StatusCreated, b)
}

// CreateDeployment godoc
// @Summary Create a new deployment
// @Description Create a new deployment for a build
// @Tags deployments
// @Accept json
// @Produce json
// @Param deployment body models.CreateDeploymentRequest true "Deployment data"
// @Success 201 {object} models.Deployment
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/deployments [post]
func (h *Handler) CreateDeployment(c *gin.Context) {
	var req models.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO deployments (build_id, environment, status, deployed_by)
		VALUES ($1, $2, 'pending', $3)
		RETURNING id, build_id, environment, status, deployed_by, deployed_at, created_at, updated_at
	`

	var d models.Deployment
	err := h.db.QueryRow(query, req.BuildID, req.Environment, req.DeployedBy).Scan(
		&d.ID, &d.BuildID, &d.Environment, &d.Status, &d.DeployedBy, &d.DeployedAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		h.log.Error("Failed to create deployment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deployment"})
		return
	}

	c.JSON(http.StatusCreated, d)
}
