package server

import (
	"context"
	"net/http"

	"github.com/nohlachilders/atlas/internal/auth"
)

// AuthenticationMiddleware is a struct that implements the http.Handler interface.
// middleware layer that tries to authenticate a user based on a given JWT.
// if successful, the userID is added to the requests context for later use under
// the key "userID" as type uuid.UUID.
type AuthenticationMiddleware struct {
	next http.Handler
	cfg  *Config
}

// AddAuthenticationMiddleware adds the AuthenticationMiddleware to the request chain.
func AddAuthenticationMiddleware(next http.Handler, cfg *Config) http.Handler {
	return &AuthenticationMiddleware{
		next: next,
		cfg:  cfg,
	}
}

// See documentation for AuthenticationMiddleware
func (h *AuthenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		h.cfg.respondWithError(w, 401, "unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, h.cfg.Secret)
	if err != nil {
		h.cfg.respondWithError(w, 401, "unauthorized", err)
		return
	}

	ctx := context.WithValue(r.Context(), MiddlewareContextKey(UserIDKey), userID)

	h.next.ServeHTTP(w, r.WithContext(ctx))
}

var UserIDKey MiddlewareContextKey = "userID"
