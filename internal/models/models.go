package models

import (
	"time"
)

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Build represents a pipeline build/run
type Build struct {
	ID          int64      `json:"id" db:"id"`
	PipelineID  int64      `json:"pipeline_id" db:"pipeline_id"`
	BuildNumber int        `json:"build_number" db:"build_number"`
	Status      string     `json:"status" db:"status"`
	Branch      string     `json:"branch" db:"branch"`
	CommitHash  string     `json:"commit_hash" db:"commit_hash"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty" db:"finished_at"`
	Duration    int        `json:"duration" db:"duration"` // in seconds
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Deployment represents a deployment
type Deployment struct {
	ID          int64     `json:"id" db:"id"`
	BuildID     int64     `json:"build_id" db:"build_id"`
	Environment string    `json:"environment" db:"environment"`
	Status      string    `json:"status" db:"status"`
	DeployedBy  string    `json:"deployed_by" db:"deployed_by"`
	DeployedAt  time.Time `json:"deployed_at" db:"deployed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// HealthCheck represents system health status
type HealthCheck struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks"`
}

// CreatePipelineRequest represents request to create a pipeline
type CreatePipelineRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdatePipelineRequest represents request to update a pipeline
type UpdatePipelineRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// CreateBuildRequest represents request to create a build
type CreateBuildRequest struct {
	PipelineID int64  `json:"pipeline_id" binding:"required"`
	Branch     string `json:"branch" binding:"required"`
	CommitHash string `json:"commit_hash" binding:"required"`
}

// CreateDeploymentRequest represents request to create a deployment
type CreateDeploymentRequest struct {
	BuildID     int64  `json:"build_id" binding:"required"`
	Environment string `json:"environment" binding:"required"`
	DeployedBy  string `json:"deployed_by" binding:"required"`
}
