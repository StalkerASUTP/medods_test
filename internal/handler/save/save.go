package save

import (
	"authentication-app/internal/util"
	"authentication-app/storage/db"
	"context"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserSaver interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}

func New(log *slog.Logger, saver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.New"
		log := log.With(slog.String("op", op))
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			log.Error("user_id is required")
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		if err := util.Validate(userID); err != nil {
			log.Error("invalid user_id", "error", err)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		tokenJWT, err := tokens.AccessGenerator(userID, )
		
	}
}

