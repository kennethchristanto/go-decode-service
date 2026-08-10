package services

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"net/url"

	"decode-service/pkg/crypto"
)

type CryptoService interface {
	Decrypt(encryptedOrder, secretKey string) (string, error)
	Encrypt(plainOrder, secretKey string) (string, error)
}

type cryptoService struct{}

func NewCryptoService() CryptoService {
	return &cryptoService{}
}

func (s *cryptoService) Decrypt(encryptedOrder, secretKey string) (string, error) {
	// URL decode twice to handle double-encoded characters
	firstDecoded, err := url.QueryUnescape(encryptedOrder)
	if err != nil {
		return "", fmt.Errorf("first URL decode failed: %w", err)
	}

	secondDecoded, err := url.QueryUnescape(firstDecoded)
	if err != nil {
		return "", fmt.Errorf("second URL decode failed: %w", err)
	}

	// Base64 decode the secret key
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("secret key Base64 decode failed: %w", err)
	}

	// Base64 decode the encrypted data
	encryptedData, err := base64.StdEncoding.DecodeString(secondDecoded)
	if err != nil {
		return "", fmt.Errorf("encrypted data Base64 decode failed: %w", err)
	}

	// Create DES cipher
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	// Check data length
	if len(encryptedData) < block.BlockSize() {
		return "", fmt.Errorf("encrypted data too short")
	}

	// ECB mode doesn't need an IV
	plaintext := make([]byte, len(encryptedData))
	ecbDecrypt(block, plaintext, encryptedData)

	// Remove PKCS5 padding
	plaintext, err = crypto.PKCS5Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to unpad data: %w", err)
	}

	return string(plaintext), nil
}

func (s *cryptoService) Encrypt(plainOrder, secretKey string) (string, error) {
	// Base64 decode the secret key
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("secret key Base64 decode failed: %w", err)
	}

	// Create DES cipher
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	// Add PKCS5 padding
	paddedData := crypto.PKCS5Pad([]byte(plainOrder), block.BlockSize())

	// ECB mode doesn't need an IV
	encryptedData := make([]byte, len(paddedData))
	ecbEncrypt(block, encryptedData, paddedData)

	// Base64 encode the encrypted data
	base64Encoded := base64.StdEncoding.EncodeToString(encryptedData)

	// URL encode twice
	firstEncoded := url.QueryEscape(base64Encoded)
	doubleEncoded := url.QueryEscape(firstEncoded)

	return doubleEncoded, nil
}

// ecbEncrypt implements ECB encryption
func ecbEncrypt(block cipher.Block, dst, src []byte) {
	if len(src)%block.BlockSize() != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		block.Encrypt(dst, src[:block.BlockSize()])
		src = src[block.BlockSize():]
		dst = dst[block.BlockSize():]
	}
}

// ecbDecrypt implements ECB decryption
func ecbDecrypt(block cipher.Block, dst, src []byte) {
	if len(src)%block.BlockSize() != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:block.BlockSize()])
		src = src[block.BlockSize():]
		dst = dst[block.BlockSize():]
	}
}