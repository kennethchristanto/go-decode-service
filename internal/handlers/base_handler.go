package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"decode-service/internal/models"
	"decode-service/internal/services"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	CryptoService      services.CryptoService
	ValidationService  services.ValidationService
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(cryptoService services.CryptoService, validationService services.ValidationService) *BaseHandler {
	return &BaseHandler{
		CryptoService:     cryptoService,
		ValidationService: validationService,
	}
}

// SendSuccess sends a successful response
func (h *BaseHandler) SendSuccess(c *gin.Context, data interface{}, message string) {
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// SendError sends an error response
func (h *BaseHandler) SendError(c *gin.Context, statusCode int, message string) {
	response := models.APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	}
	c.JSON(statusCode, response)
}

// SendValidationErrorResponse sends a validation error response
func (h *BaseHandler) SendValidationErrorResponse(c *gin.Context, message string) {
	h.SendError(c, http.StatusBadRequest, message)
}

// SendErrorResponse sends a general error response
func (h *BaseHandler) SendErrorResponse(c *gin.Context, message string) {
	h.SendError(c, http.StatusInternalServerError, message)
}