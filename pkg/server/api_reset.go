package server

import (
	"errors"
	"net/http"
)

type ResetHandler struct {
	cfg *Config
}

func (h ResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.cfg.Platform != "dev" {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "not allowed", errors.New("this should never be printed"))
		return
	}

	err := h.cfg.Database.Reset(r.Context())
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
