package server

import "net/http"

type ResetHandler struct {
	cfg *Config
}

func NewResetHandler(cfg *Config) ResetHandler {
	return ResetHandler{
		cfg: cfg,
	}
}

func (h ResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.cfg.Platform != "dev" {
		respondWithError(w, http.StatusUnauthorized, "not allowed")
		return
	}

	err := h.cfg.Database.Reset(*h.cfg.Context)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}
