package save

import (
	"authentication-app/internal/config"
	"authentication-app/internal/util/res"
	"authentication-app/internal/util/tokens"
	"authentication-app/storage/db"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserSaver interface {
	CreateUser(ctx context.Context, arg *db.CreateUserParams) (*db.User, error)
}

func New(log *slog.Logger, saver UserSaver, cfg *config.Config, tokenMan tokens.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.New"
		log := log.With(slog.String("op", op))
		id := r.URL.Query().Get("user_id")
		if id == "" {
			log.Error("user_id is required")
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		userID, err := uuid.Parse(id)
		if err != nil {
			log.Error("invalid user_id", "error", err, "user_id", userID)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		tokenJWT, err := tokenMan.AccessGenerator(userID, cfg.AccessTTL)
		if err != nil {
			log.Error("failed to generate access token", "error", err)
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := tokens.RefreshGenerator()
		if err != nil {
			log.Error("failed to generate refresh token", "error", err)
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}
		refreshTokenHash, err := tokens.RefTokenHash(refreshToken)
		if err != nil {
			log.Error("failed to hash refresh token", "error", err)
			http.Error(w, "failed to hash refresh token", http.StatusInternalServerError)
			return
		}
		_, err = saver.CreateUser(r.Context(), &db.CreateUserParams{
			ID:                    userID,
			RefreshTokenHash:      refreshTokenHash,
			RefreshTokenExpiresAt: pgtype.Timestamp{Time: time.Now().Add(cfg.RefreshTTL), Valid: true},
			UserAgent:             r.UserAgent(),
			IpAddress:             r.RemoteAddr,
			IsActive:              true,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Error("user already exists", "error", err)
				http.Error(w, "user already exists", http.StatusBadRequest)
				return

			}
			log.Error("failed to create user", "error", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		response := res.TokenResponse{
			AccessToken:  tokenJWT,
			RefreshToken: refreshToken,
		}
		res.Json(w, response, http.StatusCreated)
	}
}
