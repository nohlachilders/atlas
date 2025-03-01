package server

import (
	"net/http"

	"github.com/google/uuid"
)

type DeleteUserHandler struct {
	cfg *Config
}

func (h DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.cfg.Database.DeleteUser(r.Context(), r.Context().Value(UserIDKey).(uuid.UUID))
	if err != nil {
		h.cfg.respondWithError(w, 500, "Something went wrong", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
