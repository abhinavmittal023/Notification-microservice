package hash

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken takes length as parameter and gives secure random encoded key
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
