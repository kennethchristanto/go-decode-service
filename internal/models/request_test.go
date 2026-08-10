package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAPIResponseFormat(t *testing.T) {
	// Test that APIResponse follows the expected format
	response := APIResponse{
		Success:   true,
		Data:      map[string]string{"key": "value"},
		Message:   "Success",
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Equal(t, "Success", response.Message)
	assert.NotNil(t, response.Timestamp)
}

func TestHealthResponseFormat(t *testing.T) {
	// Test that HealthResponse follows the expected format
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	assert.Equal(t, "healthy", response.Status)
	assert.NotNil(t, response.Timestamp)
}

func TestDecryptResponseFormat(t *testing.T) {
	// Test that DecryptResponse follows the expected format
	response := DecryptResponse{
		DecryptedOrder: "A-121226-ABCDEF",
	}

	assert.Equal(t, "A-121226-ABCDEF", response.DecryptedOrder)
}

func TestEncryptResponseFormat(t *testing.T) {
	// Test that EncryptResponse follows the expected format
	response := EncryptResponse{
		EncryptedOrder: "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6",
	}

	assert.Equal(t, "qZCC%252FTvzJYUpNd9VllcaYfYKCkAKeOR6", response.EncryptedOrder)
}

func TestRequestValidation(t *testing.T) {
	// Test that request models have required fields

	// Test DecryptRequest
	decryptReq := DecryptRequest{
		EncryptedOrder: "test-encrypted",
		SecretKey:      "secretkey",
	}

	assert.NotEmpty(t, decryptReq.EncryptedOrder)
	assert.NotEmpty(t, decryptReq.SecretKey)

	// Test EncryptRequest
	encryptReq := EncryptRequest{
		PlainOrder: "A-121226-ABCDEF",
		SecretKey:  "secretkey",
	}

	assert.NotEmpty(t, encryptReq.PlainOrder)
	assert.NotEmpty(t, encryptReq.SecretKey)
}
