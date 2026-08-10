package services

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"decode-service/internal/models"
)

type ValidationService interface {
	ValidateDecryptRequest(req models.DecryptRequest) error
	ValidateEncryptRequest(req models.EncryptRequest) error
}

type ValidationServiceInternal struct{}

func NewValidationService() ValidationService {
	return &ValidationServiceInternal{}
}

// Public type for testing
type ValidationServicePublic struct {
	ValidationService
}

func (s *ValidationServiceInternal) ValidateDecryptRequest(req models.DecryptRequest) error {
	if strings.TrimSpace(req.EncryptedOrder) == "" {
		return errors.New("encrypted_order is required")
	}

	if strings.TrimSpace(req.SecretKey) == "" {
		return errors.New("secret_key is required")
	}

	// Validate Base64 format for secret key
	if !isValidBase64(req.SecretKey) {
		return errors.New("secret_key must be a valid Base64 string")
	}

	return nil
}

func (s *ValidationServiceInternal) ValidateEncryptRequest(req models.EncryptRequest) error {
	if strings.TrimSpace(req.PlainOrder) == "" {
		return errors.New("plain_order is required")
	}

	if strings.TrimSpace(req.SecretKey) == "" {
		return errors.New("secret_key is required")
	}

	// Validate Base64 format for secret key
	if !isValidBase64(req.SecretKey) {
		return errors.New("secret_key must be a valid Base64 string")
	}

	// Validate plain order format (basic validation)
	if len(req.PlainOrder) < 1 || len(req.PlainOrder) > 100 {
		return errors.New("plain_order length must be between 1 and 100 characters")
	}

	return nil
}

// IsValidBase64 checks if a string is valid Base64
func isValidBase64(s string) bool {
	// Base64 regex pattern
	base64Pattern := `^[A-Za-z0-9+/]*={0,2}$`
	match, _ := regexp.MatchString(base64Pattern, s)

	if !match {
		return false
	}

	// Check if it can be decoded
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
