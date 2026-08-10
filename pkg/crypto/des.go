package crypto

import (
	"bytes"
)

// PKCS5Pad pads data to the specified block size using PKCS5 padding
func PKCS5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS5Unpad removes PKCS5 padding from data
func PKCS5Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, nil
	}
	
	return data[:len(data)-padding], nil
}