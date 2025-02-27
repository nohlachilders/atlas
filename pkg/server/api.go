package server

import (
	"encoding/json"
	"net/http"
)

func (cfg *Config) makeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", healthReponseHandlerFunc)
	mux.Handle("POST /reset", ChainMiddlewares(
		ResetHandler{cfg: cfg},
		[]AddMiddlewareFunc{
			AddLoggingMiddleware,
		}))

	mux.Handle("POST /users", CreateUserHandler{cfg: cfg})
	mux.Handle("POST /login", LoginHandler{cfg: cfg})
	mux.Handle("POST /refresh", RefreshHandler{cfg: cfg})
	mux.Handle("POST /revoke", RevokeHandler{cfg: cfg})
}

func healthReponseHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		Error string `json:"error"`
	}
	thisError := errorStruct{
		Error: msg,
	}
	res, _ := json.Marshal(thisError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
