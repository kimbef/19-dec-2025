package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kimbef/19-dec-2025/internal/database"
	"github.com/kimbef/19-dec-2025/pkg/logger"
	"github.com/kimbef/19-dec-2025/pkg/middleware"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(db *database.DB, log *logger.Logger, isProduction bool) *gin.Engine {
	// Set Gin mode
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Create handler
	handler := NewHandler(db, log)

	// Health and readiness endpoints
	router.GET("/health", handler.HealthCheck)
	router.GET("/ready", handler.ReadyCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Pipeline routes
		pipelines := v1.Group("/pipelines")
		{
			pipelines.GET("", handler.ListPipelines)
			pipelines.POST("", handler.CreatePipeline)
			pipelines.GET("/:id", handler.GetPipeline)
			pipelines.PUT("/:id", handler.UpdatePipeline)
			pipelines.DELETE("/:id", handler.DeletePipeline)
			pipelines.GET("/:pipeline_id/builds", handler.ListBuilds)
		}

		// Build routes
		builds := v1.Group("/builds")
		{
			builds.POST("", handler.CreateBuild)
		}

		// Deployment routes
		deployments := v1.Group("/deployments")
		{
			deployments.POST("", handler.CreateDeployment)
		}
	}

	return router
}
