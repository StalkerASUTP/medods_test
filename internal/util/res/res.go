package res

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type UserIDResponse struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func Json(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
