package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AccessGenerator(userID string, ttl time.Duration) string {
	const op = "tokens.AccessGenerator"
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
	})
	return token.SignedString([]byte(conf.JWT.Secret))
}
