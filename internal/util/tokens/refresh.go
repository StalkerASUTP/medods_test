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

func RefreshGenerator() (string, error) {
	const op = "tokens.RefreshGenerator"

	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return base64.StdEncoding.EncodeToString(tokenBytes), nil
}

func RefTokenHash(token string) (string, error) {
	const op = "tokens.RefTokenHash"
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hashedBytes), nil
}
func ValidateRefToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

func GenerateHMAC(token string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}
