package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Message function takes in a message and a secret key to return the hashed version of the message
func Message(message string, key string) (string, error) {
	hash := hmac.New(sha256.New, []byte(key))
	_, err := hash.Write([]byte(message))
	if err != nil {
		return "", err
	}
	sha := hex.EncodeToString(hash.Sum(nil))
	return sha, nil
}

// Validate function takes in a message, its hash and the secret key and returns whether it is valid
func Validate(message string, messageMAC string, key string) (bool, error) {
	expectedMAC, err := Message(message, key)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(messageMAC), []byte(expectedMAC)), nil
}

// Equal function takes in two hashes and return true if they are equal
func Equal(messageMAC1 string, messageMAC2 string) bool {
	return hmac.Equal([]byte(messageMAC1), []byte(messageMAC2))
}
