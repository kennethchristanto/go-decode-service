package handlers

import (
	"github.com/gin-gonic/gin"

	"decode-service/internal/models"
	"decode-service/internal/services"
)

// EncryptHandler handles encryption requests
type EncryptHandler struct {
	*BaseHandler
}

// NewEncryptHandler creates a new encrypt handler
func NewEncryptHandler(cryptoService services.CryptoService, validationService services.ValidationService) *EncryptHandler {
	return &EncryptHandler{
		BaseHandler: NewBaseHandler(cryptoService, validationService),
	}
}

// HandleEncrypt handles the encryption endpoint
func (h *EncryptHandler) HandleEncrypt(c *gin.Context) {
	var req models.EncryptRequest
	
	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendValidationErrorResponse(c, "Invalid request format")
		return
	}
	
	// Validate request
	if err := h.ValidationService.ValidateEncryptRequest(req); err != nil {
		h.SendValidationErrorResponse(c, err.Error())
		return
	}
	
	// Perform encryption
	encrypted, err := h.CryptoService.Encrypt(req.PlainOrder, req.SecretKey)
	if err != nil {
		h.SendErrorResponse(c, "Encryption failed: "+err.Error())
		return
	}
	
	// Send success response
	h.SendSuccess(c, models.EncryptResponse{
		EncryptedOrder: encrypted,
	}, "Encryption successful")
}