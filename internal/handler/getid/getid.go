package getid

import (
	"authentication-app/internal/util/res"
	"authentication-app/internal/util/tokens"
	"authentication-app/storage/db"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*db.User, error)
}

func New(log *slog.Logger, getter UserGetter, tokenMan tokens.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.getid.New"
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
		user, err := getter.GetUserByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Error("user not found", "error", err)
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			log.Error("failed to get user by id", "error", err)
			http.Error(w, "failed to get user by id", http.StatusInternalServerError)
			return
		}
		if !user.IsActive {
			log.Error("user is not active", "error", err)
			http.Error(w, "user is not active", http.StatusUnauthorized)
			return
		}
		res.Json(w, res.UserIDResponse{UserID: user.ID}, http.StatusOK)

	}
}
