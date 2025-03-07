package server

import (
	"encoding/json"
	"net/http"

	"github.com/nohlachilders/atlas/internal/auth"
	"github.com/nohlachilders/atlas/internal/database"
)

type CreateUserHandler struct {
	cfg *Config
}

func (h CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type RequestFormat struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var data RequestFormat
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	hashed, err := auth.HashPassword(data.Password)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	err = auth.CheckPasswordHash(data.Password, hashed)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	params := database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashed,
	}
	dbUser, err := h.cfg.Database.CreateUser(r.Context(), params)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	user := User{
		Email:     dbUser.Email,
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
