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
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	refreshToken, err := h.cfg.Database.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "authentication went wrong", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "Token is expired", err)
		return
	}
	if refreshToken.RevokedAt.Valid {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "Token is revoked", err)
		return
	}

	user, err := h.cfg.Database.GetUserByID(r.Context(), refreshToken.UserID)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	timeOneHour := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(user.ID, h.cfg.Secret, timeOneHour)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	type responseShape struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, responseShape{Token: token})
}
