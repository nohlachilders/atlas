package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nohlachilders/atlas/internal/auth"
	"github.com/nohlachilders/atlas/internal/database"
)

type LoginHandler struct {
	cfg *Config
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.cfg.Database.GetUserByEmail(*h.cfg.Context, data.Email)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "email or password incorrect or something went wrong", err)
		return
	}

	err = auth.CheckPasswordHash(data.Password, user.HashedPassword)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusUnauthorized, "email or password incorrect", err)
		return
	}

	timeOneHour := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(user.ID, h.cfg.Secret, timeOneHour)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	refreshString, err := auth.MakeRefreshToken()
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	params := database.CreateRefreshTokenParams{
		Token:     refreshString,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		UserID:    user.ID,
	}
	refresh, err := h.cfg.Database.CreateRefreshToken(*h.cfg.Context, params)
	if err != nil {
		h.cfg.respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh.Token,
	})
}
