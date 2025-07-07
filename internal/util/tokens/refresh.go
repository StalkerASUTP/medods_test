package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	salt = "refresh_token_salt"
)


// RefreshGenerator generates a new refresh token
// Returns a base64 encoded token
func RefreshGenerator() (string, error) {
	const op = "tokens.RefreshGenerator"
	// Generate random bytes for the token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		// In a real application, you might want to handle this error differently
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Encode to base64
	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}

// RefTokenHash creates a bcrypt hash of the token for storage
// This function should be used before storing the token in the database
func RefTokenHash(token string) string {
	// Hash the token using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		// In a real application, you might want to handle this error differently
		panic("Failed to hash token")
	}

	return string(hashedBytes)

func ValidateRefToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// GenerateHMAC creates an HMAC signature for additional security
// This can be used for token verification if needed
func GenerateHMAC(token string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}