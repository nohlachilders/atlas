package server

import (
	"net/http"

	"github.com/google/uuid"
)

type UserInfoHandler struct {
	cfg *Config
}

func (h UserInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDContextKey).(uuid.UUID)

	dbUser, err := h.cfg.Database.GetUserByID(r.Context(), userID)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "error getting user info", err)
		return
	}

	user := User{
		Email:     dbUser.Email,
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	respondWithJSON(w, http.StatusOK, user)
}
