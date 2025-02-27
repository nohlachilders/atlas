package server

import (
	"net/http"

	"github.com/nohlachilders/atlas/internal/auth"
)

type RevokeHandler struct {
	cfg *Config
}

func (h RevokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	refreshToken, err := h.cfg.Database.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	err = h.cfg.Database.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
