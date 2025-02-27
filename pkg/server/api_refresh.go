package server

import (
	"net/http"
	"time"

	"github.com/nohlachilders/atlas/internal/auth"
)

type RefreshHandler struct {
	cfg *Config
}

func (h RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	refreshToken, err := h.cfg.Database.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authentication went wrong")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token is expired")
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token is revoked")
		return
	}

	user, err := h.cfg.Database.GetUserByID(r.Context(), refreshToken.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	timeOneHour := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(user.ID, h.cfg.Secret, timeOneHour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	type responseShape struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, responseShape{Token: token})
}
