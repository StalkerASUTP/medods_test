package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)
type TokenManager interface {
	AccessGenerator(userID uuid.UUID, ttl time.Duration) (string, error)
	AccessParser(accessToken string) (uuid.UUID, error)
}
type Manager struct {
	signingKey string
}
func NewManager(signingKey string) (*Manager,error) {
	const op = "tokens.NewManager"
	if signingKey == "" {
		return nil, fmt.Errorf("%s: signing key is required", op)
	}
	return &Manager{signingKey: signingKey}, nil
}
func (m *Manager) AccessGenerator(userID uuid.UUID, ttl time.Duration) (string, error) {
	const op = "tokens.AccessGenerator"
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	})
	accessToken, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return accessToken, nil
}

func (m *Manager) AccessParser(accessToken string) (uuid.UUID, error) {
	const op = "tokens.AccessParser"
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("%s: invalid token claims", op)
	}
	return uuid.Parse(claims["sub"].(string))
	
}