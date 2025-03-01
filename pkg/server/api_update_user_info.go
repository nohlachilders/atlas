package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nohlachilders/atlas/internal/auth"
	"github.com/nohlachilders/atlas/internal/database"
)

type UpdateUserInfoHandler struct {
	cfg *Config
}

func (h UpdateUserInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type userUpdateRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	thisRequest := userUpdateRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if thisRequest.Email == "" || thisRequest.Password == "" {
		h.cfg.respondWithError(w, http.StatusBadRequest, "Email and password are required", err)
		return
	}
	hashed, err := auth.HashPassword(thisRequest.Password)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	params := database.UpdateUserInfoParams{
		ID:             r.Context().Value(UserIDContextKey).(uuid.UUID),
		Email:          thisRequest.Email,
		HashedPassword: hashed,
	}
	user, err := h.cfg.Database.UpdateUserInfo(r.Context(), params)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
