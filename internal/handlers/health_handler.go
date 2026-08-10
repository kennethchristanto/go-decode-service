package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"decode-service/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandleHealth handles the health check endpoint
func (h *HealthHandler) HandleHealth(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}