package models

import (
	"time"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}

// DecryptRequest represents the request body for decryption
type DecryptRequest struct {
	EncryptedOrder string `json:"encrypted_order" validate:"required"`
	SecretKey     string `json:"secret_key" validate:"required"`
}

// EncryptRequest represents the request body for encryption
type EncryptRequest struct {
	PlainOrder string `json:"plain_order" validate:"required"`
	SecretKey  string `json:"secret_key" validate:"required"`
}

// DecryptResponse represents the response body for decryption
type DecryptResponse struct {
	DecryptedOrder string `json:"decrypted_order"`
}

// EncryptResponse represents the response body for encryption
type EncryptResponse struct {
	EncryptedOrder string `json:"encrypted_order"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}