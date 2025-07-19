package deactivate

import (
	"authentication-app/internal/util/res"
	"authentication-app/internal/util/tokens"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserDeactivator interface {
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
}

func New(log *slog.Logger, deactivator UserDeactivator, tokenMan tokens.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.deactivate.New"
		log := log.With(slog.String("op", op))
		token := r.URL.Query().Get("Bearer")
		if token == "" {
			log.Error("token is required")
			http.Error(w, "token is required", http.StatusBadRequest)
			return
		}
		userID, err := tokenMan.AccessParser(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Error("token expired", "error", err)
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
			log.Error("failed to parse token", "error", err)
			http.Error(w, "failed to parse token", http.StatusUnauthorized)
			return

		}

		err = deactivator.DeactivateUser(context.Background(), userID)
		if err != nil {
			log.Error("failed to deactivate user", "error", err)
			http.Error(w, "failed to deactivate user", http.StatusInternalServerError)
		}
		fmtstr := fmt.Sprintf("User: %s is diactivated", userID)
		res.Json(w, fmtstr, http.StatusOK)

	}
}
