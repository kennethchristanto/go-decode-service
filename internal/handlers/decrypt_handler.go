package handlers

import (
	"github.com/gin-gonic/gin"

	"decode-service/internal/models"
	"decode-service/internal/services"
)

// DecryptHandler handles decryption requests
type DecryptHandler struct {
	*BaseHandler
}

// NewDecryptHandler creates a new decrypt handler
func NewDecryptHandler(cryptoService services.CryptoService, validationService services.ValidationService) *DecryptHandler {
	return &DecryptHandler{
		BaseHandler: NewBaseHandler(cryptoService, validationService),
	}
}

// HandleDecrypt handles the decryption endpoint
func (h *DecryptHandler) HandleDecrypt(c *gin.Context) {
	var req models.DecryptRequest
	
	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendValidationErrorResponse(c, "Invalid request format")
		return
	}
	
	// Validate request
	if err := h.ValidationService.ValidateDecryptRequest(req); err != nil {
		h.SendValidationErrorResponse(c, err.Error())
		return
	}
	
	// Perform decryption
	decrypted, err := h.CryptoService.Decrypt(req.EncryptedOrder, req.SecretKey)
	if err != nil {
		h.SendErrorResponse(c, "Decryption failed: "+err.Error())
		return
	}
	
	// Send success response
	h.SendSuccess(c, models.DecryptResponse{
		DecryptedOrder: decrypted,
	}, "Decryption successful")
}