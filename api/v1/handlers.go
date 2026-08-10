package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"decode-service/internal/config"
	"decode-service/internal/handlers"
	"decode-service/internal/middleware"
)

// SetupRoutes sets up the API routes
func SetupRoutes(
	cfg *config.Config,
	decryptHandler *handlers.DecryptHandler,
	encryptHandler *handlers.EncryptHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	// Create Gin router
	r := gin.New()
	
	// Add middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggingMiddleware())
	
	if cfg.Security.EnableCORS {
		r.Use(middleware.CORSMiddleware())
	}
	
	// Add request size limit middleware
	if cfg.Security.MaxRequestSize > 0 {
		r.Use(gin.HandlerFunc(func(c *gin.Context) {
			if c.Request.ContentLength > cfg.Security.MaxRequestSize {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"success": false,
					"message": "Request entity too large",
					"timestamp": time.Now(),
				})
				c.Abort()
				return
			}
			c.Next()
		}))
	}
	
	// Health check endpoint
	r.GET("/health", healthHandler.HandleHealth)
	
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Decrypt endpoint
		v1.POST("/decrypt", decryptHandler.HandleDecrypt)
		
		// Encrypt endpoint
		v1.POST("/encrypt", encryptHandler.HandleEncrypt)
	}
	
	return r
}