package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncodeName 加密名称
func EncodeName(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}
